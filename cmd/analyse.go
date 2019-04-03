package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type summary struct {
	numFiles       int
	numDirectories int
	diskUsage      int64
}

type file struct {
	name string
	size int64
}

type directory struct {
	name string
}

type analysis struct {
	files       []file
	directories []directory
	diskUsage   int64
}

func humanise(bytes int64) string {
	if bytes < GB {
		if bytes < MB {
			if bytes < KB {
				return fmt.Sprintf("%i B", bytes)
			}
			return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
		}
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	}
}

func (s summary) print() {
	fmt.Println("\nSummary:")
	fmt.Println("Files:", s.numFiles)
	fmt.Println("Directories:", s.numDirectories)
	fmt.Println("Disk usage:", humanise(s.diskUsage))
}

var cmdAnalyse = &cobra.Command{
	Use:   "analyse [path]",
	Short: "Analyse directories and files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Analysing directory:", strings.Join(args, " "))
		if err := validate(cmd, args); err != nil {
			panic(err)
		}
		root := args[0]
		includeDotFiles, _ := cmd.Flags().GetBool("include-dot-files")
		summary := analyse(root, includeDotFiles)
		summary.print()
	},
}

func init() {
	var includeDotFiles bool

	rootCmd.AddCommand(cmdAnalyse)
	cmdAnalyse.Flags().BoolVarP(
		&includeDotFiles,
		"include-dot-files",
		"d",
		false,
		"include dot files (default is false)",
	)
}

func validate(cmd *cobra.Command, args []string) error {
	if err := validateCommand(cmd); err != nil {
		return err
	}
	if err := validateArgs(args); err != nil {
		return err
	}
	return nil
}

func validateCommand(cmd *cobra.Command) error {
	return nil
}

func validateArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("path not provided")
	}
	return nil
}

func isDotFile(path string) bool {
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false
}

func analyseFile(analysis *analysis, includeDotFiles bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		filename := info.Name()
		if !includeDotFiles {
			if isDotFile := isDotFile(filename); isDotFile {
				if info.IsDir() {
					fmt.Println("Skipping entire directory: " + path)
					return filepath.SkipDir
				}
				fmt.Println("Skipping file: " + path)
				return nil
			}
		}
		if info.IsDir() {
			analysis.directories = append(analysis.directories, directory{name: path})
			fmt.Println("Including file: " + path)
		} else {
			analysis.diskUsage += info.Size()
			analysis.files = append(analysis.files, file{name: path, size: info.Size()})
			fmt.Println("Including directory: " + path)
		}
		time.Sleep(50 * time.Millisecond)
		return nil
	}
}

func analyse(root string, includeDotFiles bool) summary {
	analysis := analysis{}
	err := filepath.Walk(root, analyseFile(&analysis, includeDotFiles))
	if err != nil {
		panic(err)
	}
	fmt.Println("\nFiles and directories found:")
	for _, file := range analysis.files {
		fmt.Println(file)
	}
	summary := summary{
		numFiles:       len(analysis.files),
		numDirectories: len(analysis.directories),
		diskUsage:      analysis.diskUsage,
	}
	return summary
}
