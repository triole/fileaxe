package fileaxe

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/triole/logseal"
)

type FileInfo struct {
	Path      string
	IsDir     bool
	Size      int64
	SizeHR    string
	LastMod   time.Time
	Age       time.Duration
	SortIndex string
	Err       error
}

type FileInfos []FileInfo

// Criteria describes which files Find selects and how the result is ordered.
// Ages and sizes are filters: a zero min and max disables one of them, a zero
// max leaves that range open ended.
type Criteria struct {
	Folder  string
	Matcher string
	MinAge  time.Duration
	MaxAge  time.Duration
	MinSize int64
	MaxSize int64
	SortBy  string
	Order   string
	RefTime time.Time
}

func (arr FileInfos) Len() int {
	return len(arr)
}

func (arr FileInfos) Less(i, j int) bool {
	return arr[i].SortIndex < arr[j].SortIndex
}

func (arr FileInfos) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (fa FileAxe) makeSortIndexPath(fi FileInfo) (r string) {
	r = fmt.Sprintf("%06d", len(strings.Split(fi.Path, string(os.PathSeparator))))
	r += fi.Path
	return
}

func (fa FileAxe) makeSortIndexAge(fi FileInfo) (r string) {
	return fmt.Sprintf("%050d", fi.Age.Microseconds())
}

func (fa FileAxe) Find(c Criteria) (fileList FileInfos) {
	fa.Lg.Debug("detect files", logseal.F{
		"folder":    c.Folder,
		"rxmatcher": c.Matcher,
		"min_age":   humanReadableDuration(c.MinAge),
		"max_age":   humanReadableDuration(c.MaxAge),
		"min_size":  c.MinSize,
		"max_size":  c.MaxSize,
	})
	inf, err := os.Stat(c.Folder)
	if err != nil {
		fa.Lg.IfErrFatal("can not access folder", logseal.F{
			"folder": c.Folder,
			"error":  err,
		})
		return nil
	}
	if !inf.IsDir() {
		// IfErrFatal only fires when the fields carry an error, so a plain
		// "this is not usable" has to abort through Fatal instead.
		fa.Lg.Fatal("folder to process is not a folder", logseal.F{
			"folder": c.Folder,
		})
		return nil
	}

	rxf, err := regexp.Compile(c.Matcher)
	if err != nil {
		// An unusable pattern used to reach the walk as a nil *Regexp and panic
		// there, losing the actual reason.
		fa.Lg.IfErrFatal("invalid file matcher", logseal.F{
			"regex": c.Matcher,
			"error": err,
		})
		return nil
	}
	err = filepath.Walk(c.Folder, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			fi := fa.fileInfo(path, c.RefTime)
			if fi.Err == nil && !fi.IsDir {
				fi.SortIndex = fa.makeSortIndexAge(fi)
				if c.SortBy == "path" {
					fi.SortIndex = fa.makeSortIndexPath(fi)
				}
				switch {
				case !ageInRange(fi.Age, c.MinAge, c.MaxAge):
					fa.Lg.Debug(
						"skip file, age range does not fit",
						logseal.F{"file": fi.Path, "age": humanReadableDuration(fi.Age)},
					)
				case !sizeInRange(fi.Size, c.MinSize, c.MaxSize):
					fa.Lg.Debug(
						"skip file, size range does not fit",
						logseal.F{"file": fi.Path, "size": fi.SizeHR},
					)
				default:
					fa.Lg.Debug(
						"add file",
						logseal.F{
							"file": fi.Path,
							"age":  humanReadableDuration(fi.Age),
							"size": fi.SizeHR,
						},
					)
					fileList = append(fileList, fi)
				}
			}
		}
		return nil
	})
	fa.Lg.IfErrFatal("find files failed for", logseal.F{
		"folder": c.Folder,
		"error":  err,
	})
	fa.Lg.Debug("found amount of files", logseal.F{"no": len(fileList)})

	if c.Order == "asc" {
		sort.Sort(FileInfos(fileList))
	} else {
		sort.Sort(sort.Reverse(FileInfos(fileList)))
	}
	return
}

func (fa FileAxe) fileInfo(path string, refTime time.Time) (fi FileInfo) {
	fi.Path = path
	inf, err := os.Stat(path)
	fi.Err = err
	if err != nil {
		// A dangling symlink never stats, and a file can vanish between the walk
		// and this call. Both are ordinary conditions in a live directory, so
		// report them instead of dereferencing a nil FileInfo.
		fa.Lg.IfErrInfo("retrieve file stats failed", logseal.F{
			"path": path, "error": err,
		})
		return
	}
	fi.IsDir = inf.IsDir()
	fi.Size = inf.Size()
	fi.SizeHR = humanReadableFileSize(float64(fi.Size))
	fi.LastMod = inf.ModTime()
	fi.Age = refTime.Sub(fi.LastMod)
	fa.Lg.Trace("retrieve file stats", logseal.F{
		"stats": fmt.Sprintf("%+v", fi),
	})
	return
}
