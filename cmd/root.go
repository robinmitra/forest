package cmd

import (
	"github.com/robinmitra/forest/cmd/analyse"
	"github.com/robinmitra/forest/cmd/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var VERSION string

type options struct {
	verbose bool
}

func (o *options) initialise(cmd *cobra.Command, args []string) {
	if verbose, _ := cmd.PersistentFlags().GetBool("verbose"); verbose {
		o.verbose = verbose
	}
}

func (o *options) validate() bool {
	return true
}

func (o *options) run() {
	if o.verbose {
		log.SetLevel(log.InfoLevel)
	}
}

var cmd = &cobra.Command{
	Use:   "forest",
	Short: "For the forest on your computer",
}

func NewRootCmd() *cobra.Command {
	var options = options{}

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		options.initialise(cmd, args)
		options.validate()
		options.run()
	}

	cmd.PersistentFlags().BoolVarP(&options.verbose, "verbose", "v", false, "verbose output")

	cmd.AddCommand(analyse.NewAnalyseCmd())
	cmd.AddCommand(version.NewVersionCmd(VERSION))

	return cmd
}

func Execute(version string) {
	VERSION = version

	if err := NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
