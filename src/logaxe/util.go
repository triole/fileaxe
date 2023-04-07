package logaxe

import (
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

// func nproc() int {
// 	maxProcs := runtime.GOMAXPROCS(0)
// 	numCPU := runtime.NumCPU()
// 	if maxProcs < numCPU {
// 		return maxProcs
// 	}
// 	return numCPU
// }

func humanReadableFileSize(size float64) string {
	if size < 1 {
		return "0"
	}
	var suffixes [5]string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func rxFind(rx string, content string) (r string) {
	temp, _ := regexp.Compile(rx)
	r = temp.FindString(content)
	return
}

func timestamp() string {
	dt := time.Now()
	return dt.Format("20060102t150405")
}

func str2dur(s string) (dur time.Duration, err error) {
	dur, err = str2duration.ParseDuration(s)
	return
}
