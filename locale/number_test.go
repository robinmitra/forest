package locale

import (
	"fmt"
	"testing"
)

func TestLocaliseInt(t *testing.T) {
	testCases := []struct {
		in  int64
		out string
	}{
		{in: 0, out: "0"},
		{in: 1, out: "1"},
		{in: 100, out: "100"},
		{in: 1000, out: "1,000"},
		{in: 10000, out: "10,000"},
		{in: 100000, out: "100,000"},
		{in: 1000000, out: "1,000,000"},
	}
	for _, tc := range testCases {
		// Bind the current test case as otherwise `tc` will end up referring to the last one.
		tc := tc
		t.Run(fmt.Sprintf("LocaliseInt %d", tc.in), func(t *testing.T) {
			t.Parallel()
			if res := LocaliseInt(tc.in); res != tc.out {
				t.Fatalf(
					"Expected localisation of integer %d to be %s, found %s",
					tc.in,
					tc.out,
					res,
				)
			}
		})
	}
}

func TestLocaliseFloat(t *testing.T) {
	testCases := []struct {
		in  float64
		out string
	}{
		{in: 0, out: "0.00"},
		{in: 0.0, out: "0.00"},
		{in: 0.00, out: "0.00"},
		{in: 0.001, out: "0.00"},
		{in: 0.009, out: "0.01"},
		{in: 0.01, out: "0.01"},
		{in: 0.1, out: "0.10"},
		{in: 1, out: "1.00"},
		{in: 1.2345, out: "1.23"},
		{in: 1.4567, out: "1.46"},
		{in: 100, out: "100.00"},
		{in: 1000, out: "1,000.00"},
		{in: 10000, out: "10,000.00"},
		{in: 100000, out: "100,000.00"},
		{in: 1000000, out: "1,000,000.00"},
	}
	for _, tc := range testCases {
		// Bind the current test case as otherwise `tc` will end up referring to the last one.
		tc := tc
		t.Run(fmt.Sprintf("LocaliseFloat %f", tc.in), func(t *testing.T) {
			t.Parallel()
			if res := LocaliseFloat(tc.in); res != tc.out {
				t.Fatalf(
					"Expected localisation of float %f to be %s, found %s",
					tc.in,
					tc.out,
					res,
				)
			}
		})
	}
}
