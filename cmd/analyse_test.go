package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type fileInfoMock struct {
	dir      bool
	basename string
}

func (f fileInfoMock) Name() string       { return f.basename }
func (f fileInfoMock) ModTime() time.Time { return time.Time{} }
func (f fileInfoMock) IsDir() bool        { return f.dir }
func (f fileInfoMock) Size() int64        { return int64(0) }
func (f fileInfoMock) Mode() os.FileMode {
	if f.dir {
		return 0755 | os.ModeDir
	}
	return 0644
}
func (f fileInfoMock) Sys() interface{} { return nil }

func TestInvalidArgs(t *testing.T) {
	cmd := cobra.Command{}
	var args []string
	err := validate(&cmd, args)
	if err == nil {
		t.Errorf("Expected validation to fail when passing invalid arguments.")
	}
}

func TestValidArgs(t *testing.T) {
	cmd := cobra.Command{}
	args := []string{"some-path"}
	err := validate(&cmd, args)
	if err != nil {
		t.Errorf("Expected validation to fail when passing invalid arguments.")
	}
}

func TestAnalysisOfPath(t *testing.T) {
	t.Run("Problem walking path", func(t *testing.T) {
		t.Parallel()
		analysis := newAnalysis()
		var info os.FileInfo
		var writer bytes.Buffer
		walkFunc := analyseFile(&analysis, true, &writer)
		err := walkFunc("some-path", info, errors.New("something went wrong"))
		if err == nil {
			t.Errorf("Expected error to be returned when there is a problem walking a path.")
		}
	})
	t.Run("Path is current directory", func(t *testing.T) {
		t.Parallel()
		analysis := newAnalysis()
		var info os.FileInfo
		var writer bytes.Buffer
		walkFunc := analyseFile(&analysis, true, &writer)
		err := walkFunc(".", info, nil)
		if err != nil || len(analysis.directories) > 0 || len(analysis.files) > 0 {
			t.Errorf("Expected directory to be skipped.")
		}
	})
	t.Run("Path is a dot file or directory and should be skipped", func(t *testing.T) {
		t.Parallel()
		analysis := newAnalysis()
		var writer bytes.Buffer

		// Dot directory.
		dirname := ".some-directory"
		dirInfo := fileInfoMock{dir: true, basename: dirname}
		dirWalkFunc := analyseFile(&analysis, false, &writer)
		if err := dirWalkFunc(dirname, dirInfo, nil); err != filepath.SkipDir {
			t.Errorf("Expected dot directory to be skipped.")
		}

		// Dot file.
		filename := ".some-file"
		fileInfo := fileInfoMock{dir: false, basename: filename}
		fileWalkFunc := analyseFile(&analysis, false, &writer)
		if err := fileWalkFunc(dirname, fileInfo, nil); err != nil {
			t.Errorf("Expected dot file to be skipped.")
		}

		if len(analysis.directories) > 0 || len(analysis.files) > 0 {
			t.Errorf("Expected dot files or directories to be skip.")
		}
	})
	t.Run("Path is a dot file or directory and should be included", func(t *testing.T) {
		t.Parallel()
		analysis := newAnalysis()
		var writer bytes.Buffer

		// Dot directory.
		dirname := ".some-directory"
		dirInfo := fileInfoMock{dir: true, basename: dirname}
		dirWalkFunc := analyseFile(&analysis, true, &writer)
		if err := dirWalkFunc(dirname, dirInfo, nil); err == filepath.SkipDir {
			t.Errorf("Expected dot directory to be included.")
		}

		// Dot file.
		filename := ".some-file"
		fileInfo := fileInfoMock{dir: false, basename: filename}
		fileWalkFunc := analyseFile(&analysis, true, &writer)
		if err := fileWalkFunc(dirname, fileInfo, nil); err != nil {
			t.Errorf("Expected dot file to be included.")
		}

		if len(analysis.directories) != 1 || len(analysis.files) != 1 {
			t.Errorf("Expected dot files or directories to be included.")
		}
	})
	t.Run("Path is a regular file or directory and should be included", func(t *testing.T) {
		t.Parallel()
		analysis := newAnalysis()
		var writer bytes.Buffer

		// Directory.
		dirname := "some-directory"
		dirInfo := fileInfoMock{dir: true, basename: dirname}
		dirWalkFunc := analyseFile(&analysis, false, &writer)
		if err := dirWalkFunc(dirname, dirInfo, nil); err == filepath.SkipDir {
			t.Errorf("Expected directory to be included.")
		}

		// File.
		filename := "some-file"
		fileInfo := fileInfoMock{dir: false, basename: filename}
		fileWalkFunc := analyseFile(&analysis, false, &writer)
		if err := fileWalkFunc(dirname, fileInfo, nil); err != nil {
			t.Errorf("Expected file to be included.")
		}

		if len(analysis.directories) != 1 || len(analysis.files) != 1 {
			t.Errorf("Expected files or directories to be included.")
		}
	})
}
