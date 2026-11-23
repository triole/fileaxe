package fileaxe

import (
	"os"
	"path/filepath"
	"testing"

	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

func newFileOpsTestAxe(t *testing.T) FileAxe {
	t.Helper()
	cnf := conf.InitTestConf("move", ".", `\.log$`)
	return Init(cnf, logseal.Init())
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("can not write %q: %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("can not read %q: %v", path, err)
	}
	return string(content)
}

func assertAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be absent, stat returned: %v", path, err)
	}
}

func TestMoveFileIntoExistingFolder(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "hello")

	fa := newFileOpsTestAxe(t)
	if err := fa.moveFile(FileInfo{Path: source}, targetDir); err != nil {
		t.Fatalf("moveFile failed: %v", err)
	}

	assertAbsent(t, source)
	if got := readTestFile(t, filepath.Join(targetDir, "a.log")); got != "hello" {
		t.Errorf("moved file content: got %q, want %q", got, "hello")
	}
}

// Regression test for C1: a non-existent target used to be created as a
// regular file, after which the source was removed.
func TestMoveFileMissingTargetKeepsSource(t *testing.T) {
	sourceDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "AAAA")
	target := filepath.Join(t.TempDir(), "outdir")

	fa := newFileOpsTestAxe(t)
	if err := fa.moveFile(FileInfo{Path: source}, target); err == nil {
		t.Fatal("moveFile to a non-existent target folder should fail")
	}

	if got := readTestFile(t, source); got != "AAAA" {
		t.Errorf("source was modified: got %q, want %q", got, "AAAA")
	}
	assertAbsent(t, target)
}

// Regression test for C1: a target path that is a regular file must not be
// truncated, and the source must survive.
func TestMoveFileTargetIsRegularFileKeepsSource(t *testing.T) {
	sourceDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "AAAA")
	target := filepath.Join(t.TempDir(), "outdir")
	writeTestFile(t, target, "keep me")

	fa := newFileOpsTestAxe(t)
	if err := fa.moveFile(FileInfo{Path: source}, target); err == nil {
		t.Fatal("moveFile to a regular file should fail")
	}

	if got := readTestFile(t, source); got != "AAAA" {
		t.Errorf("source was modified: got %q, want %q", got, "AAAA")
	}
	if got := readTestFile(t, target); got != "keep me" {
		t.Errorf("target file was modified: got %q, want %q", got, "keep me")
	}
}

func TestMoveFileDoesNotOverwriteDestination(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "new")
	destination := filepath.Join(targetDir, "a.log")
	writeTestFile(t, destination, "old")

	fa := newFileOpsTestAxe(t)
	if err := fa.moveFile(FileInfo{Path: source}, targetDir); err == nil {
		t.Fatal("moveFile over an existing destination should fail")
	}

	if got := readTestFile(t, source); got != "new" {
		t.Errorf("source was removed or modified: got %q, want %q", got, "new")
	}
	if got := readTestFile(t, destination); got != "old" {
		t.Errorf("destination was overwritten: got %q, want %q", got, "old")
	}
}

func TestCopyActionUsesCopyTarget(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "hello")

	cnf := conf.InitTestConf("copy", sourceDir, `\.log$`)
	cnf.Copy.Target = targetDir
	cnf.Move.Target = "" // copy must not read the move target
	fa := Init(cnf, logseal.Init())
	fa.copy(FileInfos{{Path: source}})

	if got := readTestFile(t, filepath.Join(targetDir, "a.log")); got != "hello" {
		t.Errorf("copied file content: got %q, want %q", got, "hello")
	}
	if got := readTestFile(t, source); got != "hello" {
		t.Errorf("copy removed or modified the source: got %q, want %q", got, "hello")
	}
}

func TestMoveActionUsesMoveTarget(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.log")
	writeTestFile(t, source, "hello")

	cnf := conf.InitTestConf("move", sourceDir, `\.log$`)
	cnf.Copy.Target = ""
	cnf.Move.Target = targetDir
	fa := Init(cnf, logseal.Init())
	fa.move(FileInfos{{Path: source}})

	if got := readTestFile(t, filepath.Join(targetDir, "a.log")); got != "hello" {
		t.Errorf("moved file content: got %q, want %q", got, "hello")
	}
	assertAbsent(t, source)
}
