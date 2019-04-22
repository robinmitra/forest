package analyse

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/robinmitra/forest/formatter"
)

type summary struct {
	analysis       analysis
	numFiles       int
	numDirectories int
	diskUsage      int64
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
