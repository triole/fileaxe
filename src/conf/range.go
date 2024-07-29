package conf

import (
	"strconv"
	"strings"
	"time"

	"github.com/triole/logseal"
	str2duration "github.com/xhit/go-str2duration/v2"
)

func parseDurationRangeArg(s string, lg logseal.Logseal) (min, max time.Duration) {
	var err error
	arr := strings.Split(s, ",")
	min, err = str2duration.ParseDuration(arr[0])
	lg.IfErrFatal(
		"can not parse age range arg",
		logseal.F{"string": arr[0], "error": err},
	)
	max = time.Duration(0)
	if len(arr) > 1 {
		max, err = str2duration.ParseDuration(arr[1])
		lg.IfErrFatal(
			"can not parse age range arg",
			logseal.F{"string": arr[1], "error": err},
		)
	}
	return
}

func parseNumberRangeArg(s string, lg logseal.Logseal) (min, max int) {
	var err error
	arr := strings.Split(s, ",")
	min, err = strconv.Atoi(arr[0])
	lg.IfErrFatal(
		"can not parse number range arg",
		logseal.F{"string": arr[0], "error": err},
	)
	max = min
	if len(arr) > 1 {
		max, err = strconv.Atoi(arr[1])
		lg.IfErrFatal(
			"can not parse number range arg",
			logseal.F{"string": arr[1], "error": err},
		)
	}
	return
}
