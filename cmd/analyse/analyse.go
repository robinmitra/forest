package analyse

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

type options struct {
	includeDotFiles bool
	root            string
}

func (o *options) initialise(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		o.root = args[0]
	} else {
		o.root = "."
	}
	if includeDotFiles, _ := cmd.Flags().GetBool("include-dot-files"); includeDotFiles {
		o.includeDotFiles = includeDotFiles
	}
}

func (o *options) validate() {
	if err := o.validatePath(os.Stat(o.root)); err != nil {
		log.Fatal(err)
	}
}

func (o *options) validatePath(info os.FileInfo, err error) error {
	if os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Directory \"%s\" does not exist", o.root))
	}
	return err
}

func (o *options) run() {
	log.Info("Analysing directory:", o.root)
	summary := process(o.root, o.includeDotFiles)
	summary.print()
}

var cmd = &cobra.Command{
	Use:   "analyse [path]",
	Short: "Analyse directories and files",
}

func NewAnalyseCmd() *cobra.Command {
	o := options{}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		o.initialise(cmd, args)
		o.validate()
		o.run()
	}

	cmd.Flags().BoolVarP(
		&o.includeDotFiles,
		"include-hidden-files",
		"a",
		false,
		"include hidden files (default is false)",
	)

	return cmd
}
