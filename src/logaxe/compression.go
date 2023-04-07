package logaxe

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"
	"time"

	gzip "github.com/klauspost/pgzip"
	"github.com/triole/logseal"
)

func (la LogAxe) gzipFile(sourceFile FileInfo, targetArchive string) (err error) {
	var sfil io.Reader
	var tfil io.Writer
	var tarc *gzip.Writer
	var content []byte

	start := time.Now()

	la.Lg.Info("compress file", logseal.F{
		"file": sourceFile.Path,
		"size": sourceFile.SizeHR,
	})

	if !la.DryRun {
		sfil, err = os.Open(sourceFile.Path)
		la.Lg.IfErrError(
			"can not open file",
			logseal.F{"file": sourceFile.Path, "error": err},
		)
		if err != nil {
			return
		}

		reader := bufio.NewReader(sfil)
		content, err = io.ReadAll(reader)
		la.Lg.IfErrError(
			"can not read file",
			logseal.F{"file": sourceFile.Path, "error": err},
		)
		if err != nil {
			return
		}

		tfil, err = os.Create(targetArchive)
		la.Lg.IfErrError(
			"can not init gzip writer",
			logseal.F{"file": targetArchive, "error": err},
		)
		if err != nil {
			return
		}

		tarc, err = gzip.NewWriterLevel(tfil, gzip.BestCompression)
		la.Lg.IfErrError(
			"can not write compressed archive",
			logseal.F{"file": tfil, "error": err},
		)
		if err != nil {
			return
		}

		_, err = tarc.Write(content)
		la.Lg.IfErrError("gzip error", logseal.F{"error": err})
		if err != nil {
			return
		}

		tarc.Close()

		t := time.Now()
		elapsed := t.Sub(start)

		taInfos := la.fileInfo(targetArchive, time.Now())
		la.Lg.Info(
			"compression done",
			logseal.F{
				"file": targetArchive, "duration": elapsed,
				"size": taInfos.SizeHR,
			},
		)
	}

	return
}

func (la LogAxe) makeZipArchiveFilenameAndDetectionScheme(fn string) (tar, det string) {
	folder := rxFind(".*\\/", fn)
	base := rxFind("[^/]+$", fn)
	base = rxFind(".*?\\.", base)
	base = strings.TrimSuffix(base, ".")
	tar = base + "_" + timestamp() + ".log.gz"
	det = path.Join(
		folder,
		base+"_[0-2][0-9]{3}[0-1][0-9][0-3][0-9]t[0-2][0-9][0-5][0-9][0-5][0-9]\\.log\\.gz$",
	)
	tar = path.Join(folder, tar)
	return
}
