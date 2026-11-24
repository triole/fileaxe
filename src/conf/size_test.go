package conf

import (
	"testing"

	"github.com/triole/logseal"
)

func TestParseSize(t *testing.T) {
	valid := []struct {
		in   string
		want int64
	}{
		{"0", 0},
		{"1024", 1024},
		{"1b", 1},
		{"1B", 1},
		{"1k", 1 << 10},
		{"1KB", 1 << 10},
		{"1KiB", 1 << 10},
		{"1kib", 1 << 10},
		{"1m", 1 << 20},
		{"10M", 10 << 20},
		{"1.5g", 3 << 29},
		{"1 g", 1 << 30},
		{"1gb", 1 << 30},
		{"2T", 2 << 40},
		{"1p", 1 << 50},
		{" 7k ", 7 << 10},
	}
	for _, c := range valid {
		got, err := parseSize(c.in)
		if err != nil {
			t.Errorf("parseSize(%q) failed: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseSize(%q): got %d, want %d", c.in, got, c.want)
		}
	}

	invalid := []string{"", "  ", "abc", "1x", "1.2.3g", "-1g", "1kk", "g", "."}
	for _, in := range invalid {
		if got, err := parseSize(in); err == nil {
			t.Errorf("parseSize(%q) should fail, got %d", in, got)
		}
	}
}

func TestParseSizeRangeArg(t *testing.T) {
	lg := logseal.Init()
	validateParseSizeRangeArg("1g", 1<<30, 0, lg, t)
	validateParseSizeRangeArg("1m,10m", 1<<20, 10<<20, lg, t)
	validateParseSizeRangeArg("0,0", 0, 0, lg, t)
	validateParseSizeRangeArg("0", 0, 0, lg, t)
	validateParseSizeRangeArg("512b,1k", 512, 1<<10, lg, t)
}

func validateParseSizeRangeArg(s string, expMin, expMax int64, lg logseal.Logseal, t *testing.T) {
	min, max := parseSizeRangeArg(s, lg)
	if min != expMin {
		t.Errorf("failed test parseSizeRangeArg(%q), min != expMin: %d != %d", s, min, expMin)
	}
	if max != expMax {
		t.Errorf("failed test parseSizeRangeArg(%q), max != expMax: %d != %d", s, max, expMax)
	}
}
