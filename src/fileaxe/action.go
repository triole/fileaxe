package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (fa FileAxe) list(fileList FileInfos) {
	for _, el := range fileList {
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

func (fa FileAxe) rotate(fileList FileInfos) {
	for _, fil := range fileList {
		tar := fa.makeCompressionTargetFileName(fil.Path)
		fa.Lg.Trace("make target file name",
			logseal.F{
				"source": fil.Path, "target": tar,
			},
		)

		err := fa.compressFile(fil, tar)
		if !fa.Conf.Rotate.SkipTruncate && err == nil {
			err := fa.truncateFile(fil.Path)
			fa.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil, "error": err},
			)
		} else {
			fa.Lg.Debug("skip truncate", logseal.F{"file": fil})
		}
	}
}

func (fa FileAxe) move(fileList FileInfos) {
	for _, fil := range fileList {
		fa.moveFile(fil, fa.Conf.Move.Target)
	}
}

func (fa FileAxe) remove(fileList FileInfos) {
	for _, fil := range fileList {
		if !fa.Conf.DryRun {
			if fa.Conf.Remove.Yes {
				fa.removeFile(fil.Path)
			} else {
				if askForConfirmation(fil.Path) {
					fa.removeFile(fil.Path)
				}
			}
		} else {
			fa.Lg.Info(
				"dry run, might have removed file",
				logseal.F{"path": fil.Path},
			)
		}
	}
}
