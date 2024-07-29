package conf

import (
	"strings"
	"time"

	"github.com/triole/logseal"
	str2duration "github.com/xhit/go-str2duration/v2"
)

func parseDurationRangeArg(s string, lg logseal.Logseal) (minAge, maxAge time.Duration) {
	var err error
	arr := strings.Split(s, ",")
	minAge, err = str2duration.ParseDuration(arr[0])
	lg.IfErrFatal(
		"can not parse age range arg",
		logseal.F{"string": arr[0], "error": err},
	)
	maxAge = time.Duration(0)
	if len(arr) > 1 {
		maxAge, err = str2duration.ParseDuration(arr[1])
		lg.IfErrFatal(
			"can not parse age range arg",
			logseal.F{"string": arr[1], "error": err},
		)
	}
	return
}
