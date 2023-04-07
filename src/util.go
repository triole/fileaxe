package main

import (
	"regexp"
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
