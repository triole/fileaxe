package fileaxe

import (
	"os"
	"path/filepath"
	"testing"
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
