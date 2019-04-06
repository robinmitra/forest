package disk

import (
	"testing"
)

var humaniseExamples = []struct {
	in  int64
	out string
}{
	{100, "100 B"},
	{1024, "1.00 KB"},
	{5000, "4.88 KB"},
	{1024 * 1024, "1.00 MB"},
	{5000000, "4.77 MB"},
	{1024 * 1024 * 1024, "1.00 GB"},
	{5000000000, "4.66 GB"},
}

func TestHumanise(t *testing.T) {
	for _, example := range humaniseExamples {
		res := Humanise(example.in)
		if res != example.out {
			t.Errorf(
				"Expected humanisation of %d bytes to be %s, found %s",
				example.in,
				example.out,
				res,
			)
		}
	}
}
