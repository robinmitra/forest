package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the version of Forest",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Use, version)
		},
	}
}
