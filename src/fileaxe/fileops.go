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
	Path    string
	IsDir   bool
	Size    int64
	SizeHR  string
	LastMod time.Time
	Age     time.Duration
	Err     error
}

type FileInfos []FileInfo

func (arr FileInfos) Len() int {
	return len(arr)
}

func (arr FileInfos) Less(i, j int) bool {
	si := makeSortIndex(arr[i])
	sj := makeSortIndex(arr[j])
	return si < sj
}

func (arr FileInfos) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func makeSortIndex(fi FileInfo) (r string) {
	r = fmt.Sprintf("%06d", len(strings.Split(fi.Path, string(os.PathSeparator))))
	r += fi.Path
	return
}

func (fa FileAxe) Find(basedir string, rxFilter string, maxAge time.Duration, refTime time.Time) (fileList FileInfos) {
	fa.Lg.Debug("detect files", logseal.F{
		"folder":    basedir,
		"rxmatcher": rxFilter,
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
				if maxAge == 0 {
					fileList = append(fileList, fi)
				} else {
					if fi.Age > maxAge {
						fa.Lg.Debug(
							"add file",
							logseal.F{"file": fi.Path, "age": fi.Age},
						)
						fileList = append(fileList, fi)
					} else {
						fa.Lg.Debug(
							"skip file, younger than max age",
							logseal.F{"file": fi.Path, "age": fi.Age},
						)
					}
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
	sort.Sort(FileInfos(fileList))
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

func (fa FileAxe) truncate(filename string) error {
	fa.Lg.Info("truncate", logseal.F{"file": filename})
	if !fa.Conf.DryRun {
		f, err := os.OpenFile(filename, os.O_TRUNC, 0664)
		if err != nil {
			return fmt.Errorf("could not open file %q for truncation: %v", filename, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("could not close file handler for %q after truncation: %v", filename, err)
		}
	}
	return nil
}

func (fa FileAxe) rm(filepath string) {
	if fa.Conf.DryRun {
		fa.Lg.Info(
			"dry run, would have removed file",
			logseal.F{"path": filepath},
		)
	} else {
		err := os.Remove(filepath)
		if err == nil {
			fa.Lg.Info("file removed", logseal.F{"path": filepath})
		}
		fa.Lg.IfErrError(
			"can not delete file",
			logseal.F{"path": filepath, "error": err},
		)
	}
}
