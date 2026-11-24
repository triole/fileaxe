package fileaxe

import (
	"errors"
	"fmt"

	"github.com/triole/logseal"
)

func (fa FileAxe) list(fileList FileInfos) {
	for _, fil := range fileList {
		if fa.Conf.Ls.Plain {
			fmt.Printf("%s\n", fil.Path)
		} else {
			fa.Lg.Info("list", fa.logFileInfo(fil))
		}
	}
}

func (fa FileAxe) exists(fileList FileInfos) error {
	match_no := len(fileList)
	success := fa.isInRange(match_no, fa.Conf.Exists.MinNumber, fa.Conf.Exists.MaxNumber)
	if match_no > 0 && fa.Conf.Exists.List {
		for _, fil := range fileList {
			fa.Lg.Info("list", fa.logFileInfo(fil))
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
		return fmt.Errorf(
			"exists check failed: found %d files, expected %d..%d",
			match_no, fa.Conf.Exists.MinNumber, fa.Conf.Exists.MaxNumber,
		)
	}
	return nil
}

func (fa FileAxe) compress(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		tar := fa.makeCompressionTargetFileName(fil.Path, fa.Conf.Compress.CompressionFormat)
		fa.Lg.Trace("make target file name",
			logseal.F{
				"source": fil.Path, "target": tar,
			},
		)
		if err := fa.compressFile(fil, tar, fa.Conf.Compress.CompressionFormat); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) rotate(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		tar := fa.makeCompressionTargetFileName(fil.Path, fa.Conf.Rotate.CompressionFormat)
		fa.Lg.Trace("make target file name",
			logseal.F{
				"source": fil.Path, "target": tar,
			},
		)

		if err := fa.compressFile(fil, tar, fa.Conf.Rotate.CompressionFormat); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if fa.Conf.Rotate.SkipTruncate {
			fa.Lg.Debug("skip truncate", fa.logFileInfo(fil))
			continue
		}
		if err := fa.truncateFile(fil); err != nil {
			fa.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil.Path, "size": fil.SizeHR, "error": err},
			)
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) copy(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		if err := fa.copyFile(fil, fa.Conf.Copy.Target); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) move(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		if err := fa.moveFile(fil, fa.Conf.Move.Target); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) truncate(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		if !fa.Conf.Truncate.Yes && !fa.askForConfirmation(fil.Path, "truncation") {
			continue
		}
		if err := fa.truncateFile(fil); err != nil {
			fa.Lg.IfErrError(
				"can not truncate file",
				logseal.F{"file": fil.Path, "size": fil.SizeHR, "error": err},
			)
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) remove(fileList FileInfos) (errs error) {
	for _, fil := range fileList {
		if fa.Conf.DryRun {
			fa.Lg.Info("dry run, would have removed", fa.logFileInfo(fil))
			continue
		}
		if !fa.Conf.Remove.Yes && !fa.askForConfirmation(fil.Path, "removal") {
			continue
		}
		if err := fa.removeFile(fil); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (fa FileAxe) logFileInfo(fil FileInfo) logseal.F {
	return logseal.F{
		"file": fil.Path,
		"age":  humanReadableDuration(fil.Age),
		"size": fil.SizeHR,
	}
}
