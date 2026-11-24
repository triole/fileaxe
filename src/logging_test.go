package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"fileaxe/src/conf"

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

// captureStdout runs f with os.Stdout redirected and returns what was written.
// The logger is created inside f, because logseal binds os.Stdout on Init. The
// log file has to be passed as "stdout" explicitly: without it logseal leaves
// the logrus default (stderr) in place.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("can not create pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("can not close pipe writer: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("can not read pipe: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("can not close pipe reader: %v", err)
	}
	return string(out)
}

// Regression test: the start message was logged unconditionally, so a plain
// listing printed a log line in front of its file names.
func TestLogStartRespectsPlainListing(t *testing.T) {
	cnf := conf.InitTestConf("list", "/some/folder", `\.log$`)

	cnf.Ls.Plain = true
	out := captureStdout(t, func() { logStart(logseal.Init("info", "stdout"), cnf) })
	if out != "" {
		t.Errorf("plain listing must not write a start message, got %q", out)
	}

	cnf.Ls.Plain = false
	out = captureStdout(t, func() { logStart(logseal.Init("info", "stdout"), cnf) })
	if !strings.Contains(out, "start logaxe") {
		t.Errorf("expected a start message, got %q", out)
	}
	if !strings.Contains(out, "/some/folder") {
		t.Errorf("start message should name the folder, got %q", out)
	}
}
