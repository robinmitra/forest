package cmd

import (
	"fmt"
	"github.com/robinmitra/forest/cmd/analyse"
	"github.com/spf13/cobra"
	"os"
)

var VERSION string

var rootCmd = &cobra.Command{
	Use:   "forest",
	Short: "For the forest on your computer",
}

func Execute(version string) {
	VERSION = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var verbose bool
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(analyse.NewAnalyseCmd())
}

func initConfig() {
	// TODO
}
