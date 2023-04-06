package main

import (
	"path/filepath"
	"time"

	"github.com/triole/logseal"
)

type tConf struct {
	Folder       string
	RegexMatcher string
	Now          time.Time
	MaxAgeStr    string
	MaxAge       time.Duration
	SkipTruncate bool
	DryRun       bool
}

func initConf(folder, matcher, maxAgeStr string, skipTruncate, dryRun bool) (conf tConf) {
	conf.Now = time.Now()

	abs, err := filepath.Abs(folder)
	lg.IfErrFatal("absolute file path creation failed", logseal.F{
		"error": err,
		"path":  abs,
	})
	conf.Folder = abs
	conf.RegexMatcher = matcher
	conf.SkipTruncate = skipTruncate
	conf.DryRun = dryRun

	conf.MaxAgeStr = maxAgeStr
	conf.MaxAge, err = str2dur(conf.MaxAgeStr)
	lg.IfErrFatal("can not parse max age arg", logseal.F{
		"error": err, "arg": conf.MaxAgeStr,
	})
	return
}
