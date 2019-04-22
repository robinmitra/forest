package analyse

import (
	"fmt"
	"github.com/gosuri/uilive"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func process(root string, includeDotFiles bool) summary {
	analysis := newAnalysis()
	writer := uilive.New()
	writer.RefreshInterval = time.Nanosecond
	writer.Start() // Start listening for updates and render.
	if err := filepath.Walk(root, processFile(&analysis, includeDotFiles, writer)); err != nil {
		log.Fatal(err)
	}
	if _, err := fmt.Fprintln(writer, "Done."); err != nil {
		log.Fatal(err)
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

func isDotFile(path string) bool {
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false
}

func processFile(analysis *analysis, includeDotFiles bool, w io.Writer) filepath.WalkFunc {
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
				log.Fatal(err)
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
