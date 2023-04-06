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
	LastMod time.Time
	Age     time.Duration
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
			inf, err := os.Stat(path)
			if err == nil && !inf.IsDir() {
				fi := tFileInfo{
					Path:    path,
					LastMod: inf.ModTime(),
					Age:     refTime.Sub(inf.ModTime()),
				}
				if maxAge == 0 {
					fileList = append(fileList, fi)
				} else {
					if fi.Age > maxAge {
						fileList = append(fileList, fi)
					}
				}
			} else {
				lg.IfErrInfo("stat file failed", logseal.F{
					"path": path,
				})
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
