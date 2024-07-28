package fileaxe

import (
	"fmt"
	"io"
	"os"

	"github.com/triole/logseal"
)

func (fa FileAxe) copyFile(fil FileInfo, destPath string) (err error) {
	sourcePath := fil.Path
	fa.Lg.Info(
		fa.Conf.MsgPrefix+"copy file",
		logseal.F{
			"source_age":       fil.Age,
			"source_path":      sourcePath,
			"destination_path": destPath},
	)
	if !fa.Conf.DryRun {
		var inputFile *os.File
		var outputFile *os.File
		inputFile, err = os.Open(sourcePath)
		if err != nil {
			fa.Lg.Error(
				"can not open source file",
				logseal.F{
					"source_age":       fil.Age,
					"source_path":      sourcePath,
					"destination_path": destPath},
			)
			return
		}
		defer inputFile.Close()

		outputFile, err = os.Create(destPath)
		if err != nil {
			fa.Lg.Error(
				"can not open destination file",
				logseal.F{
					"source_age":       fil.Age,
					"source_path":      sourcePath,
					"destination_path": destPath},
			)
			return
		}
		defer outputFile.Close()

		fa.Lg.Info(
			"copy file",
			logseal.F{
				"source_age":       fil.Age,
				"source_path":      sourcePath,
				"destination_path": destPath},
		)
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			fa.Lg.Error(
				"can not copy file",
				logseal.F{
					"source_age":       fil.Age,
					"source_path":      sourcePath,
					"destination_path": destPath},
			)
			return
		}
		inputFile.Close() // for Windows, close before remove
	}
	return err
}

func (fa FileAxe) moveFile(fil FileInfo, destPath string) (err error) {
	err = fa.copyFile(fil, destPath)
	if err == nil {
		fa.removeFile(fil)
	}
	return err
}

func (fa FileAxe) removeFile(fil FileInfo) (err error) {
	filepath := fil.Path
	fa.Lg.Info(fa.Conf.MsgPrefix+"remove file", logseal.F{"path": filepath})
	if !fa.Conf.DryRun {
		err = os.Remove(filepath)
		if err == nil {
			fa.Lg.Info("file removed", logseal.F{"path": filepath})
		}
		fa.Lg.IfErrError(
			"can not delete file",
			logseal.F{"path": filepath, "error": err},
		)
	}
	return
}

func (fa FileAxe) truncateFile(fil FileInfo) error {
	filepath := fil.Path
	fa.Lg.Info(fa.Conf.MsgPrefix+"truncate", logseal.F{"file": filepath})
	if !fa.Conf.DryRun {
		f, err := os.OpenFile(filepath, os.O_TRUNC, 0664)
		if err != nil {
			return fmt.Errorf("could not open file %q for truncation: %v", filepath, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("could not close file handler for %q after truncation: %v", filepath, err)
		}
	}
	return nil
}
