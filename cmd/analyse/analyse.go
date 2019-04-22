package analyse

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/gosuri/uilive"
	"github.com/robinmitra/forest/formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type summary struct {
	analysis       analysis
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

type extension struct {
	name      string
	numFiles  int
	diskUsage int64
}

type analysis struct {
	files       []file
	directories []directory
	diskUsage   int64
	extensions  map[string]extension
}

func newAnalysis() analysis {
	a := analysis{}
	a.extensions = make(map[string]extension)
	return a
}

func (a *analysis) registerExtension(extName string, size int64) {
	if len(extName) == 0 {
		extName = "(missing)"
	}
	var ext extension
	if val, ok := a.extensions[extName]; ok {
		ext = val
	}
	ext.name = extName
	ext.numFiles++
	ext.diskUsage = +size
	a.extensions[extName] = ext
}

func (a *analysis) getSortedExtensions(by string, count int) []extension {
	var extensions []extension
	for _, ext := range a.extensions {
		extensions = append(extensions, ext)
	}
	if by == "occurrence" {
		sort.Slice(extensions, func(i, j int) bool {
			return extensions[i].numFiles > extensions[j].numFiles
		})
	} else {
		sort.Slice(extensions, func(i, j int) bool {
			return extensions[i].diskUsage > extensions[j].diskUsage
		})
	}
	if count > 0 {
		if len(extensions) > count {
			return extensions[0:count]
		}
	}
	return extensions
}

func (a *analysis) getSortedFiles(by string, count int) []file {
	files := make([]file, len(a.files))
	copy(files, a.files)
	if by == "size" {
		sort.Slice(files, func(i, j int) bool {
			return files[i].size > files[j].size
		})
	}
	if count > 0 {
		if len(files) > count {
			return files[0:count]
		}
	}
	return files
}

func (s summary) print() {
	fmt.Println("\nSummary:")
	fmt.Println("\nFiles:", formatter.HumaniseNumber(int64(s.numFiles)))
	fmt.Println("Directories:", formatter.HumaniseNumber(int64(s.numDirectories)))
	fmt.Println("Disk usage:", formatter.HumaniseStorage(s.diskUsage))
	fmt.Println("")

	t := tabby.New()

	fmt.Println("Statistics:")
	fmt.Println("\nTop 5 file types by occurrence:")
	t.AddHeader("File type", "Occurrence")
	for _, ext := range s.analysis.getSortedExtensions("occurrence", 5) {
		t.AddLine(ext.name, formatter.HumaniseNumber(int64(ext.numFiles)))
	}
	t.Print()

	fmt.Println("\nTop 5 file types by total disk usage:")
	t.AddHeader("File type", "Size")
	for _, ext := range s.analysis.getSortedExtensions("size", 5) {
		t.AddLine(ext.name, formatter.HumaniseStorage(ext.diskUsage))
	}
	t.Print()

	fmt.Println("\nTop 5 files by size:")
	t.AddHeader("File", "Size")
	for _, file := range s.analysis.getSortedFiles("size", 5) {
		t.AddLine(file.name, formatter.HumaniseStorage(file.size))
	}
	t.Print()
}

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
	return
}

func (o *options) run() {
	log.Info("Analysing directory:", o.root)
	summary := analyse(o.root, o.includeDotFiles)
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

func isDotFile(path string) bool {
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false
}

func analyseFile(analysis *analysis, includeDotFiles bool, w io.Writer) filepath.WalkFunc {
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
					log.Info("Skipping entire directory: " + path)
					return filepath.SkipDir
				}
				log.Info("Skipping file: " + path)
				return nil
			}
		}
		if len(analysis.files)%5000 == 0 {
			if _, err = fmt.Fprintln(w, "Analysing..", path); err != nil {
				panic(err)
			}
		}
		if info.IsDir() {
			analysis.directories = append(analysis.directories, directory{name: filename})
			log.Info("Including file: " + path)
		} else {
			analysis.diskUsage += info.Size()
			// TODO may be create a method named registerFile which adds file and extension.
			analysis.files = append(analysis.files, file{name: filename, size: info.Size()})
			analysis.registerExtension(filepath.Ext(filename), info.Size())
			log.Info("Including directory: " + path)
		}
		return nil
	}
}

func analyse(root string, includeDotFiles bool) summary {
	analysis := newAnalysis()
	writer := uilive.New()
	writer.RefreshInterval = time.Nanosecond
	writer.Start() // Start listening for updates and render.
	if err := filepath.Walk(root, analyseFile(&analysis, includeDotFiles, writer)); err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintln(writer, "Done."); err != nil {
		panic(err)
	}
	writer.Stop()
	summary := summary{
		// TODO: Optimise this
		analysis:       analysis,
		numFiles:       len(analysis.files),
		numDirectories: len(analysis.directories),
		diskUsage:      analysis.diskUsage,
	}
	return summary
}
