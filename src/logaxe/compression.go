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

func (la LogAxe) gzipFile(sourceFile FileInfo, targetArchive string) {
	start := time.Now()

	la.Lg.Info("compress file", logseal.F{
		"file": sourceFile.Path,
		"size": sourceFile.SizeHR,
	})

	if !la.DryRun {
		sfil, _ := os.Open(sourceFile.Path)
		reader := bufio.NewReader(sfil)
		content, _ := io.ReadAll(reader)

		tfil, _ := os.Create(targetArchive)
		w, err := gzip.NewWriterLevel(tfil, gzip.BestCompression)
		la.Lg.IfErrError("can not init gzip writer", logseal.F{"error": err})

		_, err = w.Write(content)
		la.Lg.IfErrError("gzip error", logseal.F{"error": err})

		w.Close()

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
