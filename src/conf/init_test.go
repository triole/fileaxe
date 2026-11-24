package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/triole/logseal"
)

// testCLI mirrors the CLI fields that Init reads, so the wiring from flags to
// Conf can be tested without going through kong. A field added to the real CLI
// struct and read in Init has to be added here as well; without it Init panics
// on the missing key, which is the same failure mode a renamed field hits.
type testCLI struct {
	SubCommand     string
	Folder         string
	FollowSymlinks bool
	Matcher        string
	AgeRange       string
	SizeRange      string
	SortBy         string
	Order          string
	DryRun         bool
	List           struct{ Plain bool }
	Exists         struct {
		NumberRange string
		List        bool
	}
	Compress struct{ Format string }
	Rotate   struct {
		Format       string
		SkipTruncate bool
	}
	Copy     struct{ Target string }
	Move     struct{ Target string }
	Truncate struct{ Yes bool }
	Remove   struct{ Yes bool }
}

func testCLIFor(folder string) testCLI {
	cli := testCLI{
		Folder:    folder,
		Matcher:   `\.log$`,
		AgeRange:  "0,0",
		SizeRange: "0,0",
		SortBy:    "age",
		Order:     "desc",
	}
	// Init parses every range argument and aborts on a malformed one, so the
	// fixture has to carry valid defaults, not zero values.
	cli.Exists.NumberRange = "1,0"
	return cli
}

func TestInitFollowSymlinks(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real")
	if err := os.MkdirAll(real, 0755); err != nil {
		t.Fatalf("can not create %q: %v", real, err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("can not create symlink: %v", err)
	}

	cli := testCLIFor(link)
	if cnf := Init(cli, logseal.Init()); cnf.Folder != link {
		t.Errorf("by default the folder stays %q, got %q", link, cnf.Folder)
	}

	cli.FollowSymlinks = true
	cnf := Init(cli, logseal.Init())
	if cnf.Folder != real {
		t.Errorf("with --follow-symlinks the folder resolves to %q, got %q", real, cnf.Folder)
	}
	if !cnf.FollowSymlinks {
		t.Error("FollowSymlinks must be carried into Conf so it can be reported")
	}
}
