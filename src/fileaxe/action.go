package fileaxe

import (
	"github.com/triole/logseal"
)

func (la LogAxe) compress() {
	logFiles := la.Find(la.Conf.Folder, la.Conf.Matcher, 0, la.Conf.Now)
	for _, fil := range logFiles {
		tar := la.makeZipArchiveFilenameAndDetectionScheme(fil.Path)
		la.Lg.Trace("make file name and detection scheme",
			logseal.F{
				"source": fil.Path, "target": tar, "detection_scheme": tar.DetectionScheme,
			},
		)

		err := la.compressFile(fil, tar)
		if !la.Conf.SkipTruncate && err == nil {
			err := la.truncate(fil.Path)
			la.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil, "error": err},
			)
		} else {
			la.Lg.Debug("skip truncate")
		}

		if la.Conf.MaxAge.Duration > 0 {
			compressedLogs := la.Find(
				la.Conf.Folder, tar.DetectionScheme,
				la.Conf.MaxAge.Duration, la.Conf.Now,
			)
			for _, fil := range compressedLogs {
				if !la.Conf.DryRun {
					la.rm(fil.Path)
				}
			}
		}
	}
}

func (la LogAxe) remove() {
	if la.Conf.MaxAge.Duration > 0 {
		files := la.Find(
			la.Conf.Folder, la.Conf.Matcher,
			la.Conf.MaxAge.Duration, la.Conf.Now,
		)
		for _, fil := range files {
			if !la.Conf.DryRun {
				if la.Conf.Yes {
					la.rm(fil.Path)
				} else {
					if askForConfirmation(fil.Path) {
						la.rm(fil.Path)
					}
				}
			} else {
				la.Lg.Info(
					"dry run, might have removed file",
					logseal.F{"path": fil.Path},
				)
			}
		}
	} else {
		la.Lg.Info("nothing to do, remove mode requires a max age definition, use --max-age or -m")
	}
}
