package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/klauspost/pgzip"
	gzip "github.com/klauspost/pgzip"
	"github.com/triole/logseal"
)

func gzipFile(sourceFile, targetArchive string) {
	lg.Debug("compress file", logseal.F{
		"source": sourceFile,
		"target": targetArchive,
	})
	sfil, _ := os.Open(sourceFile)

	reader := bufio.NewReader(sfil)
	content, _ := ioutil.ReadAll(reader)

	tfil, _ := os.Create(targetArchive)
	w, err := gzip.NewWriterLevel(tfil, pgzip.BestCompression)
	lg.IfErrError("can not init gzip writer", logseal.F{"error": err})

	_, err = w.Write(content)
	lg.IfErrError("gzip error", logseal.F{"error": err})

	w.Close()
}

func makeZipArchiveFilenameAndDetectionScheme(fn string) (tar, det string) {
	folder := rxFind(".*\\/", fn)
	base := rxFind("[^/]+$", fn)
	base = rxFind(".*?\\.", base)
	base = strings.TrimSuffix(base, ".")
	tar = base + "_" + timestamp() + ".log.gz"
	det = "/" + base + "_[0-2][0-9]{3}[0-1][0-9][0-3][0-9]t[0-2][0-9][0-5][0-9][0-5][0-9]\\.log\\.gz$"
	tar = path.Join(folder, tar)
	return
}
