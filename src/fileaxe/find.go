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

func (fa FileAxe) Find(basedir string, rxFilter string, minAge, maxAge time.Duration, refTime time.Time) (fileList FileInfos) {
	fa.Lg.Debug("detect files", logseal.F{
		"folder":    basedir,
		"rxmatcher": rxFilter,
		"min_age":   minAge,
		"max_age":   maxAge,
	})
	inf, err := os.Stat(basedir)
	fa.Lg.IfErrFatal("can not access folder", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	if !inf.IsDir() {
		fa.Lg.IfErrFatal("not a folder", logseal.F{
			"folder": basedir,
		})
	}

	rxf, _ := regexp.Compile(rxFilter)
	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			fi := fa.fileInfo(path, refTime)
			if fi.Err == nil && !fi.IsDir {
				fi.SortIndex = fa.makeSortIndexAge(fi)
				if fa.Conf.SortBy == "path" {
					fi.SortIndex = fa.makeSortIndexPath(fi)
				}
				if fi.Age > minAge && maxAge == 0 ||
					fi.Age > minAge && fi.Age < maxAge {
					fa.Lg.Debug(
						"add file",
						logseal.F{"file": fi.Path, "age": fi.Age},
					)
					fileList = append(fileList, fi)
				} else {
					fa.Lg.Debug(
						"skip file, age range does not fit",
						logseal.F{"file": fi.Path, "age": fi.Age},
					)
				}
			}
		}
		return nil
	})
	fa.Lg.IfErrFatal("find files failed for", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	fa.Lg.Debug("found amount of files", logseal.F{"no": len(fileList)})

	if fa.Conf.Order == "asc" {
		sort.Sort(FileInfos(fileList))
	} else {
		sort.Sort(sort.Reverse(FileInfos(fileList)))
	}
	return
}

func (fa FileAxe) fileInfo(path string, refTime time.Time) (fi FileInfo) {
	inf, err := os.Stat(path)
	fi.Path = path
	fi.IsDir = inf.IsDir()
	fi.Size = inf.Size()
	fi.SizeHR = humanReadableFileSize(float64(fi.Size))
	fi.LastMod = inf.ModTime()
	fi.Age = refTime.Sub(fi.LastMod)
	fi.Err = err
	fa.Lg.Trace("retrieve file stats", logseal.F{
		"stats": fmt.Sprintf("%+v", fi), "error": err,
	})
	fa.Lg.IfErrInfo("retrieve file stats failed", logseal.F{
		"path": path, "error": err,
	})
	return
}
