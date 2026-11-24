package fileaxe

import (
	"math"
	"regexp"
	"strconv"
	"strings"
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

// humanReadableDuration renders a duration with days as the largest unit. Go's
// Duration.String() stops at hours, so a file two days old would read
// "50h0m0s". Zero components are omitted, sub-second durations keep their
// native representation where the smaller units are the useful ones.
func humanReadableDuration(d time.Duration) string {
	if d > -time.Second && d < time.Second {
		return d.String()
	}
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}

	days := d / (24 * time.Hour)
	hours := (d % (24 * time.Hour)) / time.Hour
	minutes := (d % time.Hour) / time.Minute
	seconds := (d % time.Minute) / time.Second

	var b strings.Builder
	b.WriteString(sign)
	if days > 0 {
		b.WriteString(strconv.FormatInt(int64(days), 10) + "d")
	}
	if hours > 0 {
		b.WriteString(strconv.FormatInt(int64(hours), 10) + "h")
	}
	if minutes > 0 {
		b.WriteString(strconv.FormatInt(int64(minutes), 10) + "m")
	}
	if seconds > 0 || b.Len() == len(sign) {
		b.WriteString(strconv.FormatInt(int64(seconds), 10) + "s")
	}
	return b.String()
}

func rxReplaceAllString(basestring, regex, newstring string) (r string) {
	rx := regexp.MustCompile(regex)
	r = rx.ReplaceAllString(basestring, newstring)
	return
}

func timestamp() string {
	dt := time.Now()
	return dt.Format("20060102t150405")
}
