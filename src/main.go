package main

import (
	"fmt"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()

	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	conf := initConf(
		CLI.Folder, CLI.Matcher, CLI.MaxAge, CLI.SkipTruncate, CLI.DryRun,
	)

	runProcessor(conf)
}

func runProcessor(conf tConf) {
	lg.Info("Start "+appName, logseal.F{"conf": fmt.Sprintf("%+v", conf)})
	prefix := ""
	if conf.DryRun {
		prefix = "(dry run) "
	}
	verbal := "compress and truncate"
	if conf.SkipTruncate {
		verbal = "just compress"
	}

	logFiles := find(conf.Folder, conf.RegexMatcher, 0, conf.Now)
	for _, fil := range logFiles {
		tar, det := makeZipArchiveFilenameAndDetectionScheme(fil.Path)

		lg.Trace("make file name and detection scheme",
			logseal.F{
				"source": fil.Path, "target": tar, "detection_scheme": det,
			},
		)

		lg.Info(prefix+verbal,
			logseal.F{"source": fil, "target": tar},
		)

		if !conf.DryRun {
			gzipFile(fil.Path, tar)
			if !conf.SkipTruncate {
				err := truncate(fil.Path)
				lg.IfErrError(
					"can not truncate file",
					logseal.F{"file": fil, "error": err},
				)
			} else {
				lg.Debug("skip truncate")
			}
		}

		if conf.MaxAge > 0 {
			compressedLogs := find(conf.Folder, det, conf.MaxAge, conf.Now)
			for _, fil := range compressedLogs {
				lg.Info(
					prefix+"remove file", logseal.F{
						"path": fil.Path,
						"age":  fil.Age,
					},
				)
				if !conf.DryRun {
					rm(fil.Path)
				}
			}
		}
	}
}
