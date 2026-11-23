package fileaxe

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/triole/logseal"
)

// destinationPath validates that destDir is an existing folder and returns the
// full destination path for fil. Refusing to invent a file in place of a target
// folder is what keeps a failed move from destroying the source.
func (fa FileAxe) destinationPath(destDir string, fil FileInfo) (string, error) {
	inf, err := os.Stat(destDir)
	if err != nil {
		return "", fmt.Errorf("can not access target folder %q: %w", destDir, err)
	}
	if !inf.IsDir() {
		return "", fmt.Errorf("target %q is not a folder", destDir)
	}
	return filepath.Join(destDir, filepath.Base(fil.Path)), nil
}

func (fa FileAxe) copyFile(fil FileInfo, destDir string) error {
	sourcePath := fil.Path
	destPath, err := fa.destinationPath(destDir, fil)
	if err != nil {
		fa.Lg.Error(
			fa.Conf.MsgPrefix+"can not copy file",
			logseal.F{
				"source_age":       fil.Age,
				"source_path":      sourcePath,
				"destination_path": destDir,
				"error":            err},
		)
		return err
	}
	fields := logseal.F{
		"source_age":       fil.Age,
		"source_path":      sourcePath,
		"destination_path": destPath,
	}
	fa.Lg.Info(fa.Conf.MsgPrefix+"copy file", fields)
	if fa.Conf.DryRun {
		return nil
	}

	logErr := func(err error) error {
		fa.Lg.Error(fa.Conf.MsgPrefix+"can not copy file", logseal.F{
			"source_path":      sourcePath,
			"destination_path": destPath,
			"error":            err,
		})
		return err
	}

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return logErr(fmt.Errorf("can not open source file %q: %w", sourcePath, err))
	}
	defer inputFile.Close() // closed before moveFile removes the source

	sourceInfo, err := inputFile.Stat()
	if err != nil {
		return logErr(fmt.Errorf("can not stat source file %q: %w", sourcePath, err))
	}

	// O_EXCL: never overwrite an existing destination. A move must not be able
	// to clobber a file and then delete its own source.
	outputFile, err := os.OpenFile(
		destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, sourceInfo.Mode().Perm(),
	)
	if err != nil {
		return logErr(
			fmt.Errorf("can not create destination file %q: %w", destPath, err),
		)
	}

	_, err = io.Copy(outputFile, inputFile)
	if closeErr := outputFile.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(destPath) // do not leave a partial copy behind
		return logErr(
			fmt.Errorf("can not copy %q to %q: %w", sourcePath, destPath, err),
		)
	}
	fa.Lg.Info(fa.Conf.MsgPrefix+"file copied", fields)
	return nil
}

func (fa FileAxe) moveFile(fil FileInfo, destDir string) (err error) {
	if err = fa.copyFile(fil, destDir); err != nil {
		return err
	}
	return fa.removeFile(fil)
}

func (fa FileAxe) removeFile(fil FileInfo) (err error) {
	fa.Lg.Info(fa.Conf.MsgPrefix+"remove file", fa.logFileInfo(fil))
	if !fa.Conf.DryRun {
		err = os.Remove(fil.Path)
		if err == nil {
			fa.Lg.Info("file removed", fa.logFileInfo(fil))
		}
		fa.Lg.IfErrError(
			"can not delete file",
			logseal.F{"path": fil.Path, "error": err},
		)
	}
	return
}

func (fa FileAxe) truncateFile(fil FileInfo) error {
	fa.Lg.Info(fa.Conf.MsgPrefix+"truncate", fa.logFileInfo(fil))
	if !fa.Conf.DryRun {
		f, err := os.OpenFile(fil.Path, os.O_TRUNC, 0664)
		if err != nil {
			return fmt.Errorf(
				"could not open file %q for truncation: %v", fil.Path, err,
			)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf(
				"could not close file handler for %q after truncation: %v",
				fil.Path, err,
			)
		}
	}
	return nil
}
