package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFolderArgResolvesSymlinks(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real")
	if err := os.MkdirAll(real, 0755); err != nil {
		t.Fatalf("can not create %q: %v", real, err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("can not create symlink: %v", err)
	}

	got, err := FolderArg(link, false)
	if err != nil {
		t.Fatalf("FolderArg(%q, false) failed: %v", link, err)
	}
	if got != link {
		t.Errorf(
			"without follow the argument must stay as given: got %q, want %q",
			got, link,
		)
	}

	got, err = FolderArg(link, true)
	if err != nil {
		t.Fatalf("FolderArg(%q, true) failed: %v", link, err)
	}
	if got != real {
		t.Errorf("with follow the symlink must resolve: got %q, want %q", got, real)
	}
}

func TestFolderArgIsAbsolute(t *testing.T) {
	got, err := FolderArg(".", false)
	if err != nil {
		t.Fatalf("FolderArg(\".\", false) failed: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("FolderArg must return an absolute path, got %q", got)
	}
}

func TestFolderArgMissingPath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "gone")

	// Without follow a missing folder stays a plain absolute path; the
	// existence check belongs to the walk.
	got, err := FolderArg(missing, false)
	if err != nil || got != missing {
		t.Errorf("FolderArg(%q, false) = %q, %v", missing, got, err)
	}

	// With follow, resolution has to fail instead of returning a path that
	// cannot be walked.
	if _, err := FolderArg(missing, true); err == nil {
		t.Errorf("FolderArg(%q, true) should fail for a missing path", missing)
	}
}
