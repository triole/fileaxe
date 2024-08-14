package fileaxe

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archiver/v4"
	"github.com/triole/logseal"
)

func (fa FileAxe) compressFile(sourceFile FileInfo, target, compressionFormat string) (err error) {
	start := time.Now()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{
			CompressionLevel: 9,
			Multithreaded:    true,
		},
		Archival: archiver.Tar{},
	}
	if compressionFormat == "brotli" {
		format = archiver.CompressedArchive{
			Compression: archiver.Brotli{},
			Archival:    archiver.Tar{},
		}
	}
	if compressionFormat == "bz2" {
		format = archiver.CompressedArchive{
			Compression: archiver.Bz2{},
			Archival:    archiver.Tar{},
		}
	}
	if compressionFormat == "lz4" {
		format = archiver.CompressedArchive{
			Compression: archiver.Lz4{},
			Archival:    archiver.Tar{},
		}
	}
	if compressionFormat == "snappy" {
		format = archiver.CompressedArchive{
			Compression: archiver.Sz{},
			Archival:    archiver.Tar{},
		}
	}
	if compressionFormat == "xz" {
		format = archiver.CompressedArchive{
			Compression: archiver.Xz{},
			Archival:    archiver.Tar{},
		}
	}
	fmt.Printf("%+v\n", format)
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

func (fa FileAxe) runCompression(sources []string, target string, format archiver.CompressedArchive) (err error) {
	var files []archiver.File
	fileMap := make(map[string]string)
	for _, fil := range sources {
		fileMap[fil] = filepath.Base(target)
	}

	files, err = archiver.FilesFromDisk(nil, fileMap)
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

func (fa FileAxe) makeCompressionTargetFileName(fn, compressionFormat string) (tar string) {
	cleanBase := rxReplaceAllString(filepath.Base(fn), "[^A-Za-z0-9_\\-]", "_")
	noext := strings.TrimSuffix(cleanBase, filepath.Ext(fn))

	tar = filepath.Join(
		filepath.Dir(fn), noext+"_"+timestamp()+"."+compressionFormat,
	)
	return
}
