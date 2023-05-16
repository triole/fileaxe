package logaxe

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

func (la LogAxe) compressFile(sourceFile FileInfo, target tTarget) (err error) {
	start := time.Now()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{
			CompressionLevel: 9,
			Multithreaded:    true,
		},
		Archival: archiver.Tar{},
	}
	if la.TargetFormat == "snappy" {
		format = archiver.CompressedArchive{
			Compression: archiver.Sz{},
			Archival:    archiver.Tar{},
		}
	}
	if la.TargetFormat == "xz" {
		format = archiver.CompressedArchive{
			Compression: archiver.Xz{},
			Archival:    archiver.Tar{},
		}
	}

	la.Lg.Info("compress file", logseal.F{
		"file": sourceFile.Path,
		"size": sourceFile.SizeHR,
	})

	sourceFilesArr := []string{sourceFile.Path}
	err = la.runCompression(sourceFilesArr, target, format)
	end := time.Now()
	elapsed := end.Sub(start)

	if err == nil {
		taInfos := la.fileInfo(target.FullPath, time.Now())
		la.Lg.Info(
			"compression done",
			logseal.F{
				"file": target.FullPath, "duration": elapsed,
				"size": taInfos.SizeHR,
			},
		)
	} else {
		la.Lg.Error(
			"compression failed",
			logseal.F{
				"path": sourceFile.Path, "duration": elapsed, "error": err,
			},
		)
	}
	return
}

func (la LogAxe) runCompression(sources []string, target tTarget, format archiver.CompressedArchive) (err error) {
	var files []archiver.File
	fileMap := make(map[string]string)
	for _, fil := range sources {
		fileMap[fil] = target.BaseName
	}

	files, err = archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		la.Lg.Error(
			"mapping files failed",
			logseal.F{"path": sources, "target": target, "error": err},
		)
		return err
	}

	out, err := os.Create(target.FullPath)
	if err != nil {
		la.Lg.Error(
			"can not create file",
			logseal.F{"target": target, "error": err},
		)
		return err
	}
	defer out.Close()

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		la.Lg.Error(
			"archiving files failed",
			logseal.F{"files": files, "target": target, "error": err},
		)
		return err
	}
	return nil
}

func (la LogAxe) makeZipArchiveFilenameAndDetectionScheme(fn string) (tar tTarget) {
	tar.Folder = rxFind(".*\\/", fn)
	base := rxFind("[^/]+$", fn)
	base = rxFind(".*?\\.", base)
	base = strings.TrimSuffix(base, ".")
	tar.BaseName = base + "_" + timestamp() + ".log"
	tar.ShortName = tar.BaseName + "." + la.TargetFormat
	tar.DetectionScheme = path.Join(
		tar.Folder,
		base+"_[0-2][0-9]{3}[0-1][0-9][0-3][0-9]t[0-2][0-9][0-5][0-9][0-5][0-9]\\.log\\."+la.TargetFormat+"$",
	)
	tar.FullPath = path.Join(tar.Folder, tar.ShortName)
	return
}
