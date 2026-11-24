package fileaxe

import (
	"math"
	"strings"
	"testing"
	"time"
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

func TestHumanReadableDuration(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "0s"},
		{500 * time.Millisecond, "500ms"},
		{-500 * time.Millisecond, "-500ms"},
		{1500 * time.Millisecond, "1s"},
		{90 * time.Second, "1m30s"},
		{2 * time.Hour, "2h"},
		{time.Hour + 5*time.Second, "1h5s"},
		{24 * time.Hour, "1d"},
		{36*time.Hour + 30*time.Minute, "1d12h30m"},
		{50 * time.Hour, "2d2h"},
		{49*time.Hour + 3*time.Minute + 4*time.Second, "2d1h3m4s"},
		{-90 * time.Minute, "-1h30m"},
	}
	for _, c := range cases {
		if got := humanReadableDuration(c.in); got != c.want {
			t.Errorf("humanReadableDuration(%v): got %q, want %q", c.in, got, c.want)
		}
	}
}

// Regression test: the day unit has to survive a full day, which Go's own
// duration rendering would turn into "24h0m0s".
func TestHumanReadableDurationUsesDays(t *testing.T) {
	if got := humanReadableDuration(72 * time.Hour); got != "3d" {
		t.Errorf("humanReadableDuration(72h) = %q, want %q", got, "3d")
	}
}
