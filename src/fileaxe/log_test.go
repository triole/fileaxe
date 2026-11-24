package fileaxe

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

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

// captureStdin replaces os.Stdin with a pipe carrying input and restores it
// after the test. It lets confirmation paths run without touching the terminal.
func captureStdin(t *testing.T, input string) {
	t.Helper()
	old := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("can not create stdin pipe: %v", err)
	}
	if _, err := io.WriteString(w, input); err != nil {
		t.Fatalf("can not write to stdin pipe: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("can not close stdin pipe writer: %v", err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = old
		_ = r.Close()
	})
}

func newLogTestAxe(cnf conf.Conf) FileAxe {
	return Init(cnf, logseal.Init("info", "stdout", true))
}

// Regression test: the start message was logged unconditionally, so a plain
// listing printed a log line in front of its file names.
func TestLogStartRespectsPlainListing(t *testing.T) {
	cnf := conf.InitTestConf("list", "/some/folder", `\.log$`)

	cnf.Ls.Plain = true
	out := captureStdout(t, func() { newLogTestAxe(cnf).logStart() })
	if out != "" {
		t.Errorf("plain listing must not write a start message, got %q", out)
	}

	cnf.Ls.Plain = false
	out = captureStdout(t, func() { newLogTestAxe(cnf).logStart() })
	if !strings.Contains(out, "start fileaxe") {
		t.Errorf("expected a start message, got %q", out)
	}
	for _, field := range []string{"folder=/some/folder", "rxmatcher", "symlinks=false", "min_age", "min_size"} {
		if !strings.Contains(out, field) {
			t.Errorf("start message should contain %q, got %q", field, out)
		}
	}
}

// Regression test: the start event reported the size filters as raw byte
// counts while every other size in the log was human readable.
func TestLogStartReportsHumanReadableSizeFilters(t *testing.T) {
	cnf := conf.InitTestConf("list", "/some/folder", `.*`)
	cnf.MinSize = 1 << 30
	cnf.MaxSize = 2 << 30

	out := captureStdout(t, func() { newLogTestAxe(cnf).logStart() })

	for _, field := range []string{"min_size=1GB", "max_size=2GB"} {
		if !strings.Contains(out, field) {
			t.Errorf("start message should contain %q, got %q", field, out)
		}
	}
}

// Regression test: a run used to report its start twice, once from main and
// once from Run, with different names and field vocabularies for the same
// values. Run is now the only place that reports it.
func TestRunReportsStartOnce(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.log"), "a")
	cnf := conf.InitTestConf("list", dir, `\.log$`)

	out := captureStdout(t, func() { newLogTestAxe(cnf).Run() })

	// Count any start event, not just the expected one: the original duplicate
	// reported a different name for the same moment.
	if got := strings.Count(out, "msg=start"); got != 1 {
		t.Errorf("want exactly one start message, got %d in %q", got, out)
	}
	if strings.Contains(out, "logaxe") {
		t.Errorf("start message still uses the old project name: %q", out)
	}
	if !strings.Contains(out, "a.log") {
		t.Errorf("the run should still list the matching file, got %q", out)
	}
}

// Regression test: file related info logs reported only the path and a raw
// duration, so a listing gave no size and an age past a day read "50h0m0s".
func TestLogFileInfoCarriesSizeAndReadableAge(t *testing.T) {
	fil := FileInfo{Path: "/some/file.log", SizeHR: "2KB", Age: 50 * time.Hour}
	fields := FileAxe{}.logFileInfo(fil)

	if got := fields["size"]; got != "2KB" {
		t.Errorf("size field: got %v, want %q", got, "2KB")
	}
	if got := fields["age"]; got != "2d2h" {
		t.Errorf("age field: got %v, want %q", got, "2d2h")
	}
	if got := fields["file"]; got != fil.Path {
		t.Errorf("file field: got %v, want %q", got, fil.Path)
	}
}

func TestRunLogsDoneWithStats(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "old.log")
	young := filepath.Join(dir, "young.log")
	writeTestFile(t, old, strings.Repeat("o", 1024))
	writeTestFile(t, young, strings.Repeat("y", 2048))

	now := time.Now()
	cnf := conf.InitTestConf("list", dir, `\.log$`)
	cnf.Now = now
	if err := os.Chtimes(old, now.Add(-50*time.Hour), now.Add(-50*time.Hour)); err != nil {
		t.Fatalf("can not set mtime of %q: %v", old, err)
	}
	if err := os.Chtimes(young, now.Add(-2*time.Hour), now.Add(-2*time.Hour)); err != nil {
		t.Fatalf("can not set mtime of %q: %v", young, err)
	}

	out := captureStdout(t, func() { newLogTestAxe(cnf).Run() })

	for _, field := range []string{
		"msg=done fileaxe", "action=list", "no=2", "size=3KB",
		"duration=", "oldest_age=2d2h", "newest_age=2h",
	} {
		if !strings.Contains(out, field) {
			t.Errorf("done message should contain %q, got %q", field, out)
		}
	}
	if strings.Contains(out, "errors=") {
		t.Errorf("a clean run must not report errors, got %q", out)
	}
}

func TestLogDoneOmitsAgeSpanWithoutMatches(t *testing.T) {
	cnf := conf.InitTestConf("list", t.TempDir(), `\.log$`)

	out := captureStdout(t, func() { newLogTestAxe(cnf).Run() })

	if !strings.Contains(out, "msg=done fileaxe") || !strings.Contains(out, "no=0") {
		t.Errorf("want a done message for an empty match, got %q", out)
	}
	for _, field := range []string{"oldest_age", "newest_age"} {
		if strings.Contains(out, field) {
			t.Errorf("empty match must not report %q, got %q", field, out)
		}
	}
}

func TestLogDoneRespectsPlainListing(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.log"), "a")
	cnf := conf.InitTestConf("list", dir, `\.log$`)
	cnf.Ls.Plain = true

	out := captureStdout(t, func() { _ = newLogTestAxe(cnf).Run() })

	if strings.Contains(out, "done fileaxe") {
		t.Errorf("plain listing must not write a done message, got %q", out)
	}
	if strings.TrimSpace(out) != filepath.Join(dir, "a.log") {
		t.Errorf("plain listing should only contain the file name, got %q", out)
	}
}

// Regression test: a failed action used to leave Run without an error, so the
// process exited 0 and the done event looked like a clean run.
func TestRunReportsActionFailure(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.log"), "a")
	writeTestFile(t, filepath.Join(dir, "b.log"), "b")
	cnf := conf.InitTestConf("copy", dir, `\.log$`)
	cnf.Copy.Target = filepath.Join(dir, "missing")

	var runErr error
	out := captureStdout(t, func() { runErr = newLogTestAxe(cnf).Run() })

	if runErr == nil {
		t.Fatal("Run should report that no file could be copied")
	}
	if !strings.Contains(out, "msg=done fileaxe") || !strings.Contains(out, "errors=2") {
		t.Errorf("done message should count both failures, got %q", out)
	}
	for _, name := range []string{"a.log", "b.log"} {
		if got := readTestFile(t, filepath.Join(dir, name)); got == "" {
			t.Errorf("failed copy modified %q", name)
		}
	}
}

// Regression test: a failed exists check used to call os.Exit from the action,
// which made the outcome untestable and skipped the completion event.
func TestRunReturnsErrorWhenExistsCheckFails(t *testing.T) {
	cnf := conf.InitTestConf("exists", t.TempDir(), `\.log$`)
	cnf.Exists.MinNumber = 1

	var runErr error
	out := captureStdout(t, func() { runErr = newLogTestAxe(cnf).Run() })

	if runErr == nil {
		t.Fatal("Run should report a failed exists check")
	}
	if !strings.Contains(out, "success=false") || !strings.Contains(out, "errors=1") {
		t.Errorf("want a failed check and one error in the done event, got %q", out)
	}
}

func TestCountErrors(t *testing.T) {
	one := errors.New("one")
	two := errors.New("two")
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"single", one, 1},
		{"joined", errors.Join(one, two), 2},
		{"nested", errors.Join(errors.Join(one, two), one), 3},
	}
	for _, c := range cases {
		if got := countErrors(c.err); got != c.want {
			t.Errorf("countErrors(%s): got %d, want %d", c.name, got, c.want)
		}
	}
}
