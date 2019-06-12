package browse

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
)

type options struct {
	tree bool
	root string
}

func (o *options) initialise(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		r, _ := regexp.Compile("/$")
		o.root = r.ReplaceAllString(args[0], "")
	} else {
		o.root = "."
	}
	if tree, _ := cmd.Flags().GetBool("tree"); tree {
		o.tree = tree
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
	if !o.tree {
		log.Fatal("Unknown display mode")
		return
	}
	renderTree(buildFileTree(o.root))
}

var cmd = &cobra.Command{
	Use:   "browse",
	Short: "Interactively browse directories and files",
}

func NewInteractiveCmd() *cobra.Command {
	o := options{}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		o.initialise(cmd, args)
		o.validate()
		o.run()
	}

	cmd.Flags().BoolVarP(
		&o.tree,
		"tree",
		"t",
		true,
		"browse the file tree",
	)

	return cmd
}
