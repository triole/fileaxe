package fileaxe

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mholt/archiver/v4"
	"github.com/triole/logseal"
)

type tTarget struct {
	Folder          string
	FullPath        string
	BaseName        string
	ShortName       string
	DetectionScheme string
}

func (fa FileAxe) compressFile(sourceFile FileInfo, target tTarget) (err error) {
	start := time.Now()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{
			CompressionLevel: 9,
			Multithreaded:    true,
		},
		Archival: archiver.Tar{},
	}
	if fa.Conf.Rotate.CompressionFormat == "snappy" {
		format = archiver.CompressedArchive{
			Compression: archiver.Sz{},
			Archival:    archiver.Tar{},
		}
	}
	if fa.Conf.Rotate.CompressionFormat == "xz" {
		format = archiver.CompressedArchive{
			Compression: archiver.Xz{},
			Archival:    archiver.Tar{},
		}
	}

	fa.Lg.Info("compress file", logseal.F{
		"file": sourceFile.Path,
		"size": sourceFile.SizeHR,
	})

	sourceFilesArr := []string{sourceFile.Path}
	if !fa.Conf.DryRun {
		err = fa.runCompression(sourceFilesArr, target, format)
		end := time.Now()
		elapsed := end.Sub(start)

		if err == nil {
			taInfos := fa.fileInfo(target.FullPath, time.Now())
			fa.Lg.Info(
				"compression done",
				logseal.F{
					"file": target.FullPath, "duration": elapsed,
					"size": taInfos.SizeHR,
				},
			)
		} else {
			fa.Lg.Error(
				"compression failed",
				logseal.F{
					"path": sourceFile.Path, "duration": elapsed, "error": err,
				},
			)
		}
	}
	return
}

func (fa FileAxe) runCompression(sources []string, target tTarget, format archiver.CompressedArchive) (err error) {
	var files []archiver.File
	fileMap := make(map[string]string)
	for _, fil := range sources {
		fileMap[fil] = target.BaseName
	}

	files, err = archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		fa.Lg.Error(
			"mapping files failed",
			logseal.F{"path": sources, "target": target, "error": err},
		)
		return err
	}

	out, err := os.Create(target.FullPath)
	if err != nil {
		fa.Lg.Error(
			"can not create file",
			logseal.F{"target": target, "error": err},
		)
		return err
	}
	defer out.Close()

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		fa.Lg.Error(
			"archiving files failed",
			logseal.F{"files": files, "target": target, "error": err},
		)
		return err
	}
	return nil
}

func (fa FileAxe) makeZipArchiveFilenameAndDetectionScheme(fn string) (tar tTarget) {
	tar.Folder = rxFind(".*\\/", fn)
	base := rxFind("[^/]+$", fn)
	base = rxFind(".*?\\.", base)
	base = strings.TrimSuffix(base, ".")
	tar.BaseName = base + "_" + timestamp() + ".log"
	tar.ShortName = tar.BaseName + "." + fa.Conf.Rotate.CompressionFormat
	tar.DetectionScheme = path.Join(
		tar.Folder,
		base+"_[0-2][0-9]{3}[0-1][0-9][0-3][0-9]t[0-2][0-9][0-5][0-9][0-5][0-9]\\.log\\."+fa.Conf.Rotate.CompressionFormat+"$",
	)
	tar.FullPath = path.Join(tar.Folder, tar.ShortName)
	return
}
