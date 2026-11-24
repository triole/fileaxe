package fileaxe

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

// Regression test: a failed archive used to leave a truncated .gz behind. The
// next run could mistake that file for a valid archive, and the source it was
// supposed to preserve had no usable counterpart.
func TestCompressFileRemovesBrokenArchive(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "unreadable.log")
	writeTestFile(t, source, "payload")
	if err := os.Chmod(source, 0); err != nil {
		t.Fatalf("can not remove read permission from %q: %v", source, err)
	}
	t.Cleanup(func() { _ = os.Chmod(source, 0o644) })

	// A privileged runner ignores the mode, so the failure cannot be provoked.
	if f, err := os.Open(source); err == nil {
		f.Close()
		t.Skip("source stays readable, probably running as a privileged user")
	}

	fa := newFileOpsTestAxe(t)
	target := filepath.Join(dir, "unreadable.gz")
	fil := FileInfo{Path: source, Size: 7, SizeHR: "7B"}

	if err := fa.compressFile(fil, target, "gz"); err == nil {
		t.Fatal("compressing an unreadable file should fail")
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("failed compression left archive %q behind (stat: %v)", target, err)
	}

	if err := os.Chmod(source, 0o644); err != nil {
		t.Fatalf("can not restore read permission on %q: %v", source, err)
	}
	if got := readTestFile(t, source); got != "payload" {
		t.Errorf("source content: got %q, want %q", got, "payload")
	}
}

// Regression test: a.b.log and a_b.log used to sanitize to the same target, so
// one archive silently replaced the other.
func TestCompressionTargetNamesStayDistinct(t *testing.T) {
	dir := t.TempDir()
	fa := FileAxe{}

	dotted := fa.makeCompressionTargetFileName(filepath.Join(dir, "a.b.log"), "gz")
	underscored := fa.makeCompressionTargetFileName(filepath.Join(dir, "a_b.log"), "gz")

	if dotted == underscored {
		t.Errorf("distinct sources share the target %q", dotted)
	}
	for _, target := range []string{dotted, underscored} {
		if filepath.Dir(target) != dir {
			t.Errorf("target %q is not in the source folder %q", target, dir)
		}
		if filepath.Ext(target) != ".gz" {
			t.Errorf("target %q does not carry the compression suffix", target)
		}
	}
}

// Regression test: an existing archive was truncated by os.Create, destroying
// the only copy of an earlier compression.
func TestCompressRefusesToOverwriteExistingArchive(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "a.log")
	writeTestFile(t, source, "payload")

	fa := newFileOpsTestAxe(t)
	target := fa.makeCompressionTargetFileName(source, "gz")
	writeTestFile(t, target, "sentinel")

	err := fa.compressFile(FileInfo{Path: source, Size: 7, SizeHR: "7B"}, target, "gz")
	if err == nil {
		t.Fatal("compressing onto an existing archive should fail")
	}
	if got := readTestFile(t, target); got != "sentinel" {
		t.Errorf("existing archive was modified: got %q, want %q", got, "sentinel")
	}
	if got := readTestFile(t, source); got != "payload" {
		t.Errorf("source content: got %q, want %q", got, "payload")
	}
}

// Regression test: rotate truncated both sources after writing their colliding
// target, so one file's content was lost while both sources were emptied.
func TestRotateKeepsBothCollidingFiles(t *testing.T) {
	dir := t.TempDir()
	sources := map[string]string{
		"a.b.log": "dotted content",
		"a_b.log": "underscored content",
	}
	for name, content := range sources {
		writeTestFile(t, filepath.Join(dir, name), content)
	}

	cnf := conf.InitTestConf("rotate", dir, `\.log$`)
	cnf.Rotate.CompressionFormat = "gz"
	fa := Init(cnf, logseal.Init())
	if err := fa.Run(); err != nil {
		t.Fatalf("rotate failed: %v", err)
	}

	archives := filesWithSuffix(t, dir, ".gz")
	if len(archives) != 2 {
		t.Fatalf("want one archive per source, got %v", archives)
	}

	seen := map[string]string{}
	for _, archive := range archives {
		for name, content := range tarGzContents(t, archive) {
			if _, dup := seen[name]; dup {
				t.Errorf("entry %q is archived more than once", name)
			}
			seen[name] = content
		}
	}
	for name, want := range sources {
		if got, ok := seen[name]; !ok || got != want {
			t.Errorf("archive entry %q: got %q (present=%v), want %q", name, got, ok, want)
		}
		fi, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("can not stat %q: %v", name, err)
		}
		if fi.Size() != 0 {
			t.Errorf("source %q should be truncated, size is %d", name, fi.Size())
		}
	}
}

func tarGzContents(t *testing.T, archive string) map[string]string {
	t.Helper()
	f, err := os.Open(archive)
	if err != nil {
		t.Fatalf("can not open %q: %v", archive, err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("can not read gzip stream of %q: %v", archive, err)
	}
	defer gz.Close()

	contents := map[string]string{}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("can not read tar stream of %q: %v", archive, err)
		}
		content, err := io.ReadAll(tr)
		if err != nil {
			t.Fatalf("can not read %q from %q: %v", hdr.Name, archive, err)
		}
		contents[filepath.Base(hdr.Name)] = string(content)
	}
	return contents
}
