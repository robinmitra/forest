package analyse

import "sort"

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

func newAnalysis() analysis {
	a := analysis{}
	a.extensions = make(map[string]extension)
	return a
}
