package cmd

import (
	"github.com/robinmitra/forest/cmd/analyse"
	"github.com/robinmitra/forest/cmd/browse"
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
	var o = options{}

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		o.initialise(cmd, args)
		o.validate()
		o.run()
	}

	cmd.PersistentFlags().BoolVarP(&o.verbose, "verbose", "v", false, "verbose output")

	cmd.AddCommand(analyse.NewAnalyseCmd())
	cmd.AddCommand(version.NewVersionCmd(VERSION))
	cmd.AddCommand(browse.NewInteractiveCmd())

	return cmd
}

func Execute(version string) {
	VERSION = version

	if err := NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
