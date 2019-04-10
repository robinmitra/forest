package formatter

import (
	"fmt"
	"testing"
)

func TestHumanise(t *testing.T) {
	testCases := []struct {
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
	for _, tc := range testCases {
		// Bind the current test case as otherwise `tc` will end up referring to the last one.
		tc := tc
		t.Run(fmt.Sprintf("HumaniseStorage %d bytes", tc.in), func(t *testing.T) {
			t.Parallel()
			if res := HumaniseStorage(tc.in); res != tc.out {
				t.Fatalf(
					"Expected humanisation of %d bytes to be %s, found %s",
					tc.in,
					tc.out,
					res,
				)
			}
		})
	}
}
