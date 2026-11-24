package fileaxe

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

// readTarGzEntry returns the content of the entry called name inside a gzip
// compressed tar archive, which is what the compression commands produce.
func readTarGzEntry(t *testing.T, archive, name string) string {
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

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("can not read tar stream of %q: %v", archive, err)
		}
		if filepath.Base(hdr.Name) != name {
			continue
		}
		content, err := io.ReadAll(tr)
		if err != nil {
			t.Fatalf("can not read %q from %q: %v", name, archive, err)
		}
		return string(content)
	}
	t.Fatalf("archive %q does not contain %q", archive, name)
	return ""
}

func filesWithSuffix(t *testing.T, dir, suffix string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("can not read %q: %v", dir, err)
	}
	var res []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			res = append(res, filepath.Join(dir, entry.Name()))
		}
	}
	return res
}

func assertFileMagic(t *testing.T, path string, magic []byte) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("can not open %q: %v", path, err)
	}
	defer f.Close()
	head := make([]byte, len(magic))
	if _, err := io.ReadFull(f, head); err != nil {
		t.Fatalf("can not read %d bytes of %q: %v", len(magic), path, err)
	}
	if !bytes.Equal(head, magic) {
		t.Errorf("%q starts with % x, want % x", path, head, magic)
	}
}

// Regression test: rotate used to read the compression format of the compress
// command, so -g/--format was silently ignored and every archive was gz.
func TestRotateUsesRotateCompressionFormat(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "a.log")
	writeSizedFile(t, source, 1024)

	cnf := conf.InitTestConf("rotate", dir, `\.log$`)
	cnf.Rotate.CompressionFormat = "xz"
	cnf.Compress.CompressionFormat = "gz" // must not leak into rotate
	fa := Init(cnf, logseal.Init())
	fa.Run()

	archives := filesWithSuffix(t, dir, ".xz")
	if len(archives) != 1 {
		t.Fatalf("want exactly one .xz archive, got %v", archives)
	}
	if gz := filesWithSuffix(t, dir, ".gz"); len(gz) != 0 {
		t.Errorf("rotate wrote gz archives although xz was requested: %v", gz)
	}
	assertFileMagic(t, archives[0], []byte{0xfd, '7', 'z', 'X', 'Z', 0x00})

	if fi, err := os.Stat(source); err != nil {
		t.Fatalf("can not stat %q: %v", source, err)
	} else if fi.Size() != 0 {
		t.Errorf("source should be truncated after rotation, size is %d", fi.Size())
	}
}
