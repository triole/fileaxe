package fileaxe

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archives"
	"github.com/triole/logseal"
)

func (fa FileAxe) compressFile(sourceFile FileInfo, target, compressionFormat string) (err error) {
	start := time.Now()

	format := archives.CompressedArchive{
		Compression: archives.Gz{
			CompressionLevel: 9,
			Multithreaded:    true,
		},
		Archival: archives.Tar{},
	}
	if compressionFormat == "brotli" {
		format = archives.CompressedArchive{
			Compression: archives.Brotli{},
			Archival:    archives.Tar{},
		}
	}
	if compressionFormat == "bz2" {
		format = archives.CompressedArchive{
			Compression: archives.Bz2{},
			Archival:    archives.Tar{},
		}
	}
	if compressionFormat == "lz4" {
		format = archives.CompressedArchive{
			Compression: archives.Lz4{},
			Archival:    archives.Tar{},
		}
	}
	if compressionFormat == "snappy" {
		format = archives.CompressedArchive{
			Compression: archives.Sz{},
			Archival:    archives.Tar{},
		}
	}
	if compressionFormat == "xz" {
		format = archives.CompressedArchive{
			Compression: archives.Xz{},
			Archival:    archives.Tar{},
		}
	}
	fa.Lg.Info(fa.Conf.MsgPrefix+"compress file", logseal.F{
		"file":   sourceFile.Path,
		"size":   sourceFile.SizeHR,
		"target": target,
	})

	sourceFilesArr := []string{sourceFile.Path}
	if !fa.Conf.DryRun {
		err = fa.runCompression(sourceFilesArr, target, format)
		end := time.Now()
		elapsed := end.Sub(start)

		if err == nil {
			taInfos := fa.fileInfo(target, time.Now())
			fa.Lg.Info(
				"compression done",
				logseal.F{
					"file": target, "duration": elapsed,
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

func (fa FileAxe) runCompression(sources []string, target string, format archives.CompressedArchive) (err error) {
	var files []archives.FileInfo
	ctx := context.TODO()
	fileMap := make(map[string]string)
	for _, fil := range sources {
		fileMap[fil] = filepath.Base(fil)
	}

	files, err = archives.FilesFromDisk(ctx, nil, fileMap)
	if err != nil {
		fa.Lg.Error(
			"mapping files failed",
			logseal.F{"path": sources, "target": target, "error": err},
		)
		return err
	}

	out, err := os.Create(target)
	if err != nil {
		fa.Lg.Error(
			"can not create file",
			logseal.F{"target": target, "error": err},
		)
		return err
	}
	defer out.Close()

	err = format.Archive(ctx, out, files)
	if err != nil {
		fa.Lg.Error(
			"archiving files failed",
			logseal.F{"files": files, "target": target, "error": err},
		)
		return err
	}
	return nil
}

func (fa FileAxe) makeCompressionTargetFileName(fn, compressionFormat string) (tar string) {
	cleanBase := rxReplaceAllString(filepath.Base(fn), "[^A-Za-z0-9_\\-]", "_")
	noext := strings.TrimSuffix(cleanBase, filepath.Ext(fn))

	tar = filepath.Join(
		filepath.Dir(fn), noext+"_"+timestamp()+"."+compressionFormat,
	)
	return
}
