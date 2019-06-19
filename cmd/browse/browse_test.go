package browse

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func TestInvalidPath(t *testing.T) {
	cmd := cobra.Command{}
	var args []string
	o := options{}
	o.initialise(&cmd, args)
	var info os.FileInfo

	err := o.validatePath(info, errors.New("something went wrong"))

	if err == nil {
		t.Errorf("Expected validation to fail when passing invalid path.")
	}
}

func TestValidPath(t *testing.T) {
	cmd := cobra.Command{}
	var args []string
	o := options{}
	o.initialise(&cmd, args)
	var info os.FileInfo

	err := o.validatePath(info, nil)

	if err != nil {
		t.Errorf("Expected validation to pass when passing valid path.")
	}
}
