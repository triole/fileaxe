package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/triole/logseal"
)

func TestLogTarget(t *testing.T) {
	cases := []struct {
		file   string
		dryRun bool
		want   string
	}{
		{"/var/log/fileaxe.log", true, "stdout"},
		{"/dev/stdout", true, "stdout"},
		{"", true, "stdout"},
		{"/var/log/fileaxe.log", false, "/var/log/fileaxe.log"},
		{"/dev/stdout", false, "/dev/stdout"},
		{"", false, ""},
	}
	for _, c := range cases {
		if got := logTarget(c.file, c.dryRun); got != c.want {
			t.Errorf(
				"logTarget(%q, dryRun=%v): got %q, want %q",
				c.file, c.dryRun, got, c.want,
			)
		}
	}
}

// Regression test: a plain listing wrote log lines to stdout in front of its
// file names, because the logger always went to stdout by default.
func TestRouteLogsForPlainListing(t *testing.T) {
	for _, target := range []string{"", "stdout", "/dev/stdout"} {
		lg := logseal.Init("info", target)
		routeLogsForPlainListing(lg, target, true)
		if lg.Logrus.Out != os.Stderr {
			t.Errorf("target %q: plain listing should log to stderr", target)
		}

		lg = logseal.Init("info", target)
		routeLogsForPlainListing(lg, target, false)
		if lg.Logrus.Out != os.Stdout {
			t.Errorf("target %q: a normal run should keep logging to stdout", target)
		}
	}
}

func TestPlainListingKeepsLogFile(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "fileaxe.log")
	lg := logseal.Init("info", logFile)

	routeLogsForPlainListing(lg, logFile, true)
	lg.Info("this line belongs in the file")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("can not read %q: %v", logFile, err)
	}
	if !strings.Contains(string(content), "this line belongs in the file") {
		t.Errorf("an explicit log file must keep receiving logs, got %q", content)
	}
}
