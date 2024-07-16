package logaxe

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	return arr[i].Path < arr[j].Path
}

func (arr FileInfos) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (la LogAxe) Find(basedir string, rxFilter string, maxAge time.Duration, refTime time.Time) (fileList FileInfos) {
	la.Lg.Debug("detect files", logseal.F{
		"folder":    basedir,
		"rxmatcher": rxFilter,
		"max_age":   maxAge,
	})
	inf, err := os.Stat(basedir)
	la.Lg.IfErrFatal("can not access folder", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	if !inf.IsDir() {
		la.Lg.IfErrFatal("not a folder", logseal.F{
			"folder": basedir,
		})
	}

	rxf, _ := regexp.Compile(rxFilter)
	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			fi := la.fileInfo(path, refTime)
			if fi.Err == nil && !fi.IsDir {
				if maxAge == 0 {
					fileList = append(fileList, fi)
				} else {
					if fi.Age > maxAge {
						la.Lg.Debug(
							"add file",
							logseal.F{"file": fi.Path, "age": fi.Age},
						)
						fileList = append(fileList, fi)
					} else {
						la.Lg.Debug(
							"skip file, younger than max age",
							logseal.F{"file": fi.Path, "age": fi.Age},
						)
					}
				}
			}
		}
		return nil
	})
	la.Lg.IfErrFatal("find files failed for", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	la.Lg.Debug("found amount of files", logseal.F{"no": len(fileList)})
	sort.Sort(FileInfos(fileList))
	return
}

func (la LogAxe) fileInfo(path string, refTime time.Time) (fi FileInfo) {
	inf, err := os.Stat(path)
	fi.Path = path
	fi.IsDir = inf.IsDir()
	fi.Size = inf.Size()
	fi.SizeHR = humanReadableFileSize(float64(fi.Size))
	fi.LastMod = inf.ModTime()
	fi.Age = refTime.Sub(fi.LastMod)
	fi.Err = err
	la.Lg.Trace("retrieve file stats", logseal.F{
		"stats": fmt.Sprintf("%+v", fi), "error": err,
	})
	la.Lg.IfErrInfo("retrieve file stats failed", logseal.F{
		"path": path, "error": err,
	})
	return
}

func (la LogAxe) truncate(filename string) error {
	la.Lg.Info("truncate", logseal.F{"file": filename})
	if !la.Conf.DryRun {
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

func (la LogAxe) rm(filepath string) {
	err := os.Remove(filepath)
	la.Lg.IfErrError(
		"can not delete file",
		logseal.F{"path": filepath, "error": err},
	)
}
