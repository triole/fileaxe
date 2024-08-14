package fileaxe

import (
	"fmt"
	"os"

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

func (fa FileAxe) exists(fileList FileInfos) {
	match_no := len(fileList)
	success := fa.isInRange(match_no, fa.Conf.Exists.MinNumber, fa.Conf.Exists.MaxNumber)
	if match_no > 0 && fa.Conf.Exists.List {
		for _, el := range fileList {
			fa.Lg.Info(
				el.Path,
				logseal.F{"age": el.Age, "lastmod": el.LastMod},
			)
		}
	}
	fa.Lg.Info(
		"exists check results",
		logseal.F{
			"exp_min": fa.Conf.Exists.MinNumber,
			"exp_max": fa.Conf.Exists.MaxNumber,
			"no":      match_no,
			"success": success,
		},
	)
	if !success {
		os.Exit(1)
	}
}

func (fa FileAxe) compress(fileList FileInfos) {
	for _, fil := range fileList {
		tar := fa.makeCompressionTargetFileName(fil.Path)
		fa.Lg.Trace("make target file name",
			logseal.F{
				"source": fil.Path, "target": tar,
			},
		)
		err := fa.compressFile(fil, tar)
		fa.Lg.IfErrError(
			"can not truncate file",
			logseal.F{"file": fil, "error": err},
		)
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
			err := fa.truncateFile(fil)
			fa.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil, "error": err},
			)
		} else {
			fa.Lg.Debug("skip truncate", logseal.F{"file": fil})
		}
	}
}

func (fa FileAxe) copy(fileList FileInfos) {
	for _, fil := range fileList {
		fa.copyFile(fil, fa.Conf.Move.Target)
	}
}

func (fa FileAxe) move(fileList FileInfos) {
	for _, fil := range fileList {
		fa.moveFile(fil, fa.Conf.Move.Target)
	}
}

func (fa FileAxe) truncate(fileList FileInfos) {
	for _, fil := range fileList {
		if fa.Conf.Truncate.Yes {
			fa.truncateFile(fil)
		} else {
			if fa.askForConfirmation(fil.Path, "truncation") {
				fa.truncateFile(fil)
			}
		}
	}
}

func (fa FileAxe) remove(fileList FileInfos) {
	for _, fil := range fileList {
		if !fa.Conf.DryRun {
			if fa.Conf.Remove.Yes {
				fa.removeFile(fil)
			} else {
				if fa.askForConfirmation(fil.Path, "removal") {
					fa.removeFile(fil)
				}
			}
		}
	}
}
