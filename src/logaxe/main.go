package logaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (la LogAxe) Run() {
	la.Lg.Info("Start logaxe", logseal.F{"conf": fmt.Sprintf("%+v", la)})

	if la.DryRun {
		la.Lg.Info(" --- DRY RUN START ---")
	}

	logFiles := la.Find(la.Folder, la.RegexMatcher, 0, la.Now)
	for _, fil := range logFiles {
		tar, det := la.makeZipArchiveFilenameAndDetectionScheme(fil.Path)
		la.Lg.Trace("make file name and detection scheme",
			logseal.F{
				"source": fil.Path, "target": tar, "detection_scheme": det,
			},
		)

		la.gzipFile(fil, tar)
		if !la.SkipTruncate {
			err := la.truncate(fil.Path)
			la.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil, "error": err},
			)
		} else {
			la.Lg.Debug("skip truncate")
		}

		if la.MaxAge > 0 {
			compressedLogs := la.Find(la.Folder, det, la.MaxAge, la.Now)
			for _, fil := range compressedLogs {
				la.Lg.Info(
					"remove file", logseal.F{
						"path": fil.Path,
						"age":  fil.Age,
					},
				)
				if !la.DryRun {
					la.rm(fil.Path)
				}
			}
		}
	}

	if la.DryRun {
		la.Lg.Info(" --- DRY RUN END ---")
	}
}
