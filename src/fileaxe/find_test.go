package fileaxe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

func newFindTestAxe(t *testing.T) FileAxe {
	t.Helper()
	return Init(conf.InitTestConf("list", ".", `.*`), logseal.Init())
}

func writeSizedFile(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.WriteFile(path, []byte(strings.Repeat("x", size)), 0644); err != nil {
		t.Fatalf("can not write %q: %v", path, err)
	}
}

func TestFindFiltersBySize(t *testing.T) {
	dir := t.TempDir()
	for name, size := range map[string]int{
		"small.log": 10,
		"one.log":   1 << 20,
		"two.log":   2 << 20,
		"five.log":  5 << 20,
	} {
		writeSizedFile(t, filepath.Join(dir, name), size)
	}

	fa := newFindTestAxe(t)
	cases := []struct {
		name     string
		min, max int64
		want     int
	}{
		{"no filter", 0, 0, 4},
		{"min only, inclusive", 1 << 20, 0, 3},
		{"bounded, both inclusive", 1 << 20, 2 << 20, 2},
		{"nothing big enough", 6 << 20, 0, 0},
		{"exact single size", 10, 10, 1},
	}
	for _, c := range cases {
		got, err := fa.Find(Criteria{
			Folder:  dir,
			Matcher: `\.log$`,
			MinSize: c.min,
			MaxSize: c.max,
			RefTime: time.Now(),
		})
		if err != nil {
			t.Fatalf("%s: find failed: %v", c.name, err)
		}
		if len(got) != c.want {
			t.Errorf("%s: got %d files, want %d", c.name, len(got), c.want)
		}
		for _, fil := range got {
			if fil.Size < c.min && c.min != 0 {
				t.Errorf("%s: file %q with size %d is below min %d", c.name, fil.Path, fil.Size, c.min)
			}
			if c.max != 0 && fil.Size > c.max {
				t.Errorf("%s: file %q with size %d is above max %d", c.name, fil.Path, fil.Size, c.max)
			}
		}
	}
}

func TestRotateRespectsSizeFilter(t *testing.T) {
	dir := t.TempDir()
	big := filepath.Join(dir, "big.log")
	small := filepath.Join(dir, "small.log")
	bigContent := strings.Repeat("B", 2<<20)
	writeTestFile(t, big, bigContent)
	writeTestFile(t, small, strings.Repeat("s", 1024))

	cnf := conf.InitTestConf("rotate", dir, `\.log$`)
	cnf.MinSize = 1 << 20
	fa := Init(cnf, logseal.Init())
	if err := fa.Run(); err != nil {
		t.Fatalf("rotate failed: %v", err)
	}

	// The source is only truncated after a successful compression, so a missing
	// or empty archive would mean the rotated content was lost.
	archives := filesWithSuffix(t, dir, ".gz")
	if len(archives) != 1 {
		t.Fatalf("want exactly one archive, got %v", archives)
	}
	if !strings.HasPrefix(filepath.Base(archives[0]), "big_log_") {
		t.Errorf(
			"archive %q does not belong to the rotated file %q",
			filepath.Base(archives[0]), filepath.Base(big),
		)
	}
	if fi, err := os.Stat(archives[0]); err != nil {
		t.Fatalf("can not stat %q: %v", archives[0], err)
	} else if fi.Size() == 0 {
		t.Errorf("archive %q is empty", archives[0])
	}
	if got := readTarGzEntry(t, archives[0], "big.log"); got != bigContent {
		t.Errorf(
			"archive %q holds %d bytes of the original content, want %d",
			filepath.Base(archives[0]), len(got), len(bigContent),
		)
	}

	if fi, err := os.Stat(big); err != nil {
		t.Fatalf("can not stat %q: %v", big, err)
	} else if fi.Size() != 0 {
		t.Errorf("big file should be rotated and truncated, size is %d", fi.Size())
	}
	if fi, err := os.Stat(small); err != nil {
		t.Fatalf("can not stat %q: %v", small, err)
	} else if fi.Size() != 1024 {
		t.Errorf("small file is below the size filter, size is %d, want 1024", fi.Size())
	}
}

// Regression test: a stat failure used to be stored in FileInfo.Err only after
// FileInfo had already been built from the nil FileInfo, which panicked.
func TestFileInfoReportsStatFailure(t *testing.T) {
	fa := newFindTestAxe(t)
	missing := filepath.Join(t.TempDir(), "missing.log")

	fi := fa.fileInfo(missing, time.Now())

	if fi.Err == nil {
		t.Fatalf("want an error for %q, got none", missing)
	}
	if fi.Path != missing {
		t.Errorf("path is %q, want %q", fi.Path, missing)
	}
	if fi.IsDir || fi.Size != 0 {
		t.Errorf("a failed stat must not report file details, got dir=%v size=%d", fi.IsDir, fi.Size)
	}
}

// Regression test: a dangling symlink made Find panic instead of skipping the
// unreadable entry.
func TestFindSkipsUnstatableFiles(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.log"), "a")
	broken := filepath.Join(dir, "broken.log")
	if err := os.Symlink(filepath.Join(dir, "gone.log"), broken); err != nil {
		t.Skipf("can not create symlink: %v", err)
	}

	fa := newFindTestAxe(t)
	files, err := fa.Find(Criteria{Folder: dir, Matcher: `\.log$`, RefTime: time.Now()})
	if err != nil {
		t.Fatalf("an unstatable entry must not fail the walk: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("want only the readable file, got %d files: %v", len(files), files)
	}
	if filepath.Base(files[0].Path) != "a.log" {
		t.Errorf("got %q, want %q", filepath.Base(files[0].Path), "a.log")
	}
}

// Regression test: an unreadable directory used to disappear from the result
// without an error, so a partial sweep still looked successful.
func TestFindReportsUnreadableDirectory(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "visible.log"), "a")
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o755); err != nil {
		t.Fatalf("can not create %q: %v", locked, err)
	}
	writeTestFile(t, filepath.Join(locked, "hidden.log"), "b")
	if err := os.Chmod(locked, 0); err != nil {
		t.Fatalf("can not lock %q: %v", locked, err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	// A privileged runner can read the directory anyway, so the error cannot be
	// provoked there.
	if _, err := os.ReadDir(locked); err == nil {
		t.Skip("locked directory stays readable, probably running as a privileged user")
	}

	fa := newFindTestAxe(t)
	files, err := fa.Find(Criteria{Folder: dir, Matcher: `\.log$`, RefTime: time.Now()})
	if err == nil {
		t.Fatal("an unreadable directory should be reported")
	}
	if len(files) != 1 || filepath.Base(files[0].Path) != "visible.log" {
		t.Errorf("the readable file should still be listed, got %v", files)
	}
}

// runFindCapture runs Find with a logger that records fatal calls instead of
// exiting the test process, and returns what was selected, how often the run
// aborted, what was logged and the find error.
func runFindCapture(t *testing.T, cnf conf.Conf, c Criteria) (FileInfos, int, string, error) {
	t.Helper()
	var exits int
	var files FileInfos
	var err error
	out := captureStdout(t, func() {
		lg := logseal.Init("info", "stdout", true)
		// logrus exits on Fatal; keep the test process alive and count instead.
		lg.Logrus.ExitFunc = func(int) { exits++ }
		files, err = Init(cnf, lg).Find(c)
	})
	return files, exits, out, err
}

// Regression test: an unusable matcher reached the walk as a nil *Regexp and
// panicked instead of reporting the pattern.
func TestFindRejectsInvalidMatcher(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.log"), "a")
	cnf := conf.InitTestConf("list", dir, "[")

	files, exits, out, findErr := runFindCapture(t, cnf, Criteria{
		Folder: dir, Matcher: "[", RefTime: time.Now(),
	})

	if findErr == nil {
		t.Error("an invalid matcher should be returned as an error")
	}

	if len(files) != 0 {
		t.Errorf("an invalid matcher must not select files, got %d", len(files))
	}
	if exits != 1 {
		t.Errorf("want the run to abort exactly once, got %d exits", exits)
	}
	if !strings.Contains(out, "invalid file matcher") || !strings.Contains(out, "regex=[") {
		t.Errorf("want the rejected pattern in the log, got %q", out)
	}
}

// Regression test: a regular file passed as folder was accepted, because the
// guard used IfErrFatal without an error field and therefore never fired.
func TestFindRejectsNonFolderPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.log")
	writeTestFile(t, path, "a")
	cnf := conf.InitTestConf("list", path, `.*`)

	files, exits, out, findErr := runFindCapture(t, cnf, Criteria{
		Folder: path, Matcher: `.*`, RefTime: time.Now(),
	})

	if findErr == nil {
		t.Error("a regular file as folder should be returned as an error")
	}

	if len(files) != 0 {
		t.Errorf("a regular file must not be processed as folder, got %d files", len(files))
	}
	if exits != 1 {
		t.Errorf("want the run to abort exactly once, got %d exits", exits)
	}
	if !strings.Contains(out, "is not a folder") || !strings.Contains(out, path) {
		t.Errorf("want the rejected path in the log, got %q", out)
	}
}

// Regression test: a missing folder used to be dereferenced instead of
// reported, because the stat error was only logged.
func TestFindReportsMissingFolder(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "gone")
	cnf := conf.InitTestConf("list", missing, `.*`)

	files, exits, out, findErr := runFindCapture(t, cnf, Criteria{
		Folder: missing, Matcher: `.*`, RefTime: time.Now(),
	})

	if findErr == nil {
		t.Error("a missing folder should be returned as an error")
	}

	if len(files) != 0 {
		t.Errorf("a missing folder must not produce files, got %d", len(files))
	}
	if exits != 1 {
		t.Errorf("want the run to abort exactly once, got %d exits", exits)
	}
	if !strings.Contains(out, "can not access folder") || !strings.Contains(out, missing) {
		t.Errorf("want the missing path in the log, got %q", out)
	}
}

// Regression test: a symlinked folder was walked as a symlink, so nothing was
// found and nothing was reported. With --follow-symlinks conf resolves the
// argument and the walk sees the target.
func TestFindFollowsResolvedFolder(t *testing.T) {
	real := t.TempDir()
	writeTestFile(t, filepath.Join(real, "a.log"), "a")
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("can not create symlink: %v", err)
	}

	fa := newFindTestAxe(t)
	criteria := Criteria{Folder: link, Matcher: `\.log$`, RefTime: time.Now()}

	got, err := fa.Find(criteria)
	if err != nil {
		t.Fatalf("find of a symlinked root failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("without resolution a symlinked root is not walked, got %d files", len(got))
	}

	resolved, err := conf.FolderArg(link, true)
	if err != nil {
		t.Fatalf("can not resolve %q: %v", link, err)
	}
	criteria.Folder = resolved

	got, err = fa.Find(criteria)
	if err != nil {
		t.Fatalf("find of the resolved folder failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("the resolved folder should list a.log, got %d files", len(got))
	}
	if filepath.Base(got[0].Path) != "a.log" {
		t.Errorf("got %q, want %q", filepath.Base(got[0].Path), "a.log")
	}
}
