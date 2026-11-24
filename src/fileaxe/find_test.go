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
		got := fa.Find(Criteria{
			Folder:  dir,
			Matcher: `\.log$`,
			MinSize: c.min,
			MaxSize: c.max,
			RefTime: time.Now(),
		})
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
	fa.Run()

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
	files := fa.Find(Criteria{Folder: dir, Matcher: `\.log$`, RefTime: time.Now()})

	if len(files) != 1 {
		t.Fatalf("want only the readable file, got %d files: %v", len(files), files)
	}
	if filepath.Base(files[0].Path) != "a.log" {
		t.Errorf("got %q, want %q", filepath.Base(files[0].Path), "a.log")
	}
}
