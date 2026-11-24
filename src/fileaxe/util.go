package fileaxe

import (
	"math"
	"regexp"
	"strconv"
	"time"
)

func humanReadableFileSize(size float64) string {
	if size < 1 {
		return "0"
	}
	idx := 0
	val := size
	for val >= 1024 && idx < len(sizeSuffixes)-1 {
		val /= 1024
		idx++
	}
	val = math.Round(val*100) / 100
	return strconv.FormatFloat(val, 'f', -1, 64) + sizeSuffixes[idx]
}

// sizeSuffixes are the units humanReadableFileSize reports, in binary steps. An
// int64 tops out at 8 EB, but the loop stops at the last suffix instead of
// indexing past it, so larger float values still render.
var sizeSuffixes = [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func rxReplaceAllString(basestring, regex, newstring string) (r string) {
	rx := regexp.MustCompile(regex)
	r = rx.ReplaceAllString(basestring, newstring)
	return
}

func timestamp() string {
	dt := time.Now()
	return dt.Format("20060102t150405")
}
