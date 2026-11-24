package fileaxe

import (
	"math"
	"strings"
	"testing"
)

func TestHumanReadableFileSize(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{0.5, "0"},
		{1, "1B"},
		{1023, "1023B"},
		{1024, "1KB"},
		{1536, "1.5KB"},
		{1 << 20, "1MB"},
		{1 << 30, "1GB"},
		{1 << 40, "1TB"},
		{1 << 50, "1PB"},
		{1 << 60, "1EB"},
		{math.MaxInt64, "8EB"},
	}
	for _, c := range cases {
		if got := humanReadableFileSize(c.in); got != c.want {
			t.Errorf("humanReadableFileSize(%v): got %q, want %q", c.in, got, c.want)
		}
	}
}

// Regression test: sizes of 1 PiB and above panicked with an index out of
// range, because the suffix table stopped at TB.
func TestHumanReadableFileSizeIsTotal(t *testing.T) {
	for _, size := range []float64{1 << 60, math.MaxInt64, 1e30} {
		got := humanReadableFileSize(size)
		if !strings.HasSuffix(got, "EB") {
			t.Errorf(
				"humanReadableFileSize(%v) = %q, want an EB suffix",
				size, got,
			)
		}
	}
}
