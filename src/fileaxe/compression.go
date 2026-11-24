package fileaxe

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
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
					"file":        target,
					"duration":    humanReadableDuration(elapsed),
					"source_size": sourceFile.SizeHR,
					"size":        taInfos.SizeHR,
				},
			)
		} else {
			fa.Lg.Error(
				"compression failed",
				logseal.F{
					"path":        sourceFile.Path,
					"source_size": sourceFile.SizeHR,
					"duration":    humanReadableDuration(elapsed),
					"error":       err,
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

	// O_EXCL: never overwrite an existing archive. Two sources can sanitize to
	// the same base name, and a rerun within the same second can repeat a target;
	// truncating an earlier archive would destroy the only copy of its content.
	out, err := os.OpenFile(
		target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666,
	)
	if err != nil {
		fa.Lg.Error(
			"can not create archive",
			logseal.F{"target": target, "error": err},
		)
		return err
	}

	archiveErr := format.Archive(ctx, out, files)
	// Close before removing: a failed archive must not be left behind, and an
	// open file can not be removed on every platform. A close error counts as a
	// failed archive, because the tail of the archive may never have reached disk.
	if closeErr := out.Close(); archiveErr == nil {
		archiveErr = closeErr
	}
	if archiveErr != nil {
		if removeErr := os.Remove(target); removeErr != nil {
			fa.Lg.Error(
				"can not remove broken archive",
				logseal.F{"target": target, "error": removeErr},
			)
		}
		fa.Lg.Error(
			"archiving files failed",
			logseal.F{"files": files, "target": target, "error": archiveErr},
		)
		return archiveErr
	}
	return nil
}

func (fa FileAxe) makeCompressionTargetFileName(fn, compressionFormat string) (tar string) {
	cleanBase := rxReplaceAllString(filepath.Base(fn), "[^A-Za-z0-9_\\-]", "_")

	tar = filepath.Join(
		filepath.Dir(fn),
		cleanBase+"_"+timestamp()+"_"+sourceTag(fn)+"."+compressionFormat,
	)
	return
}

// sourceTag returns a short stable tag for the full source path. Sanitizing a
// base name maps different files onto the same string (a.b.log and a_b.log both
// become a_b_log), so the tag keeps their archives apart even when they share a
// directory and a timestamp.
func sourceTag(path string) string {
	sum := sha256.Sum256([]byte(path))
	return hex.EncodeToString(sum[:6])
}
