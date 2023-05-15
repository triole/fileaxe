package logaxe

import (
	"path/filepath"
	"time"

	"github.com/triole/logseal"
)

type LogAxe struct {
	Folder       string
	RegexMatcher string
	TargetFormat string
	Now          time.Time
	MaxAgeStr    string
	MaxAge       time.Duration
	SkipTruncate bool
	DryRun       bool
	Lg           logseal.Logseal
}

func InitLogAxe(folder, matcher, format, maxAgeStr string, skipTruncate, dryRun bool, lg logseal.Logseal) (la LogAxe) {
	la.Now = time.Now()
	la.Lg = lg

	abs, err := filepath.Abs(folder)
	la.Lg.IfErrFatal("absolute file path creation failed", logseal.F{
		"error": err,
		"path":  abs,
	})
	la.Folder = abs
	la.RegexMatcher = matcher
	la.TargetFormat = format
	la.SkipTruncate = skipTruncate
	la.DryRun = dryRun

	la.MaxAgeStr = maxAgeStr
	la.MaxAge, err = str2dur(la.MaxAgeStr)
	la.Lg.IfErrFatal("can not parse max age arg", logseal.F{
		"error": err, "arg": la.MaxAgeStr,
	})
	return
}
