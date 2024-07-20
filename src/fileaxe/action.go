package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (fa FileAxe) list() {
	logFiles := fa.Find(fa.Conf.Folder, fa.Conf.Matcher, 0, fa.Conf.Now)
	for _, el := range logFiles {
		if fa.Conf.Ls.Plain {
			fmt.Printf("%s\n", el.Path)
		} else {
			fa.Lg.Info(
				el.Path,
				logseal.F{"age": el.Age, "lastmod": el.LastMod},
			)
		}
	}
}

func (fa FileAxe) rotate() {
	logFiles := fa.Find(fa.Conf.Folder, fa.Conf.Matcher, 0, fa.Conf.Now)
	for _, fil := range logFiles {
		tar := fa.makeZipArchiveFilenameAndDetectionScheme(fil.Path)
		fa.Lg.Trace("make file name and detection scheme",
			logseal.F{
				"source": fil.Path, "target": tar, "detection_scheme": tar.DetectionScheme,
			},
		)

		err := fa.compressFile(fil, tar)
		if !fa.Conf.Rotate.SkipTruncate && err == nil {
			err := fa.truncate(fil.Path)
			fa.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil, "error": err},
			)
		} else {
			fa.Lg.Debug("skip truncate")
		}

		if fa.Conf.MaxAge > 0 {
			compressedLogs := fa.Find(
				fa.Conf.Folder, tar.DetectionScheme,
				fa.Conf.MaxAge, fa.Conf.Now,
			)
			for _, fil := range compressedLogs {
				if !fa.Conf.DryRun {
					fa.rm(fil.Path)
				}
			}
		}
	}
}

func (fa FileAxe) remove() {
	if fa.Conf.MaxAge > 0 {
		files := fa.Find(
			fa.Conf.Folder, fa.Conf.Matcher,
			fa.Conf.MaxAge, fa.Conf.Now,
		)
		for _, fil := range files {
			if !fa.Conf.DryRun {
				if fa.Conf.Remove.Yes {
					fa.rm(fil.Path)
				} else {
					if askForConfirmation(fil.Path) {
						fa.rm(fil.Path)
					}
				}
			} else {
				fa.Lg.Info(
					"dry run, might have removed file",
					logseal.F{"path": fil.Path},
				)
			}
		}
	} else {
		fa.Lg.Info("nothing to do, remove mode requires a max age definition, use --max-age or -m")
	}
}
