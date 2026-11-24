package fileaxe

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
