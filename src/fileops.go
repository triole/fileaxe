package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/triole/logseal"
)

type tFileInfo struct {
	Path    string
	IsDir   bool
	LastMod time.Time
	Age     time.Duration
	Err     error
}

type tFileInfos []tFileInfo

func (arr tFileInfos) Len() int {
	return len(arr)
}

func (arr tFileInfos) Less(i, j int) bool {
	return arr[i].Path < arr[j].Path
}

func (arr tFileInfos) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func find(basedir string, rxFilter string, maxAge time.Duration, refTime time.Time) (fileList tFileInfos) {
	lg.Debug("detect files", logseal.F{
		"folder":    basedir,
		"rxmatcher": rxFilter,
		"max_age":   maxAge,
	})
	inf, err := os.Stat(basedir)
	lg.IfErrFatal("can not access folder", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	if !inf.IsDir() {
		lg.IfErrFatal("not a folder", logseal.F{
			"folder": basedir,
		})
	}

	rxf, _ := regexp.Compile(rxFilter)
	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			fi := fileInfo(path, refTime)
			if fi.Err == nil && !fi.IsDir {
				if maxAge == 0 {
					fileList = append(fileList, fi)
				} else {
					if fi.Age > maxAge {
						fileList = append(fileList, fi)
					}
				}
			}
		}
		return nil
	})
	lg.IfErrFatal("find files failed for", logseal.F{
		"folder": basedir,
		"error":  err,
	})
	lg.Debug("found amount of files", logseal.F{"no": len(fileList)})
	sort.Sort(tFileInfos(fileList))
	return
}

func fileInfo(path string, refTime time.Time) (fi tFileInfo) {
	inf, err := os.Stat(path)
	fi.Path = path
	fi.IsDir = inf.IsDir()
	fi.LastMod = inf.ModTime()
	fi.Age = refTime.Sub(fi.LastMod)
	fi.Err = err
	lg.IfErrInfo("retrieve file stats failed", logseal.F{
		"path": path, "error": err,
	})
	return
}

func truncate(filename string) error {
	lg.Debug("truncate file", logseal.F{"file": filename})
	f, err := os.OpenFile(filename, os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("could not open file %q for truncation: %v", filename, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close file handler for %q after truncation: %v", filename, err)
	}
	return nil
}

func rm(filepath string) {
	err := os.Remove(filepath)
	lg.IfErrError(
		"can not delete file",
		logseal.F{"path": filepath, "error": err},
	)
}
