package conf

import (
	"path/filepath"
	"time"

	"github.com/triole/logseal"
	str2duration "github.com/xhit/go-str2duration/v2"
)

type Conf struct {
	Now          time.Time
	Folder       string
	Matcher      string
	TargetFormat string
	MaxAge       tMaxAge
	Remove       bool
	Yes          bool
	SkipTruncate bool
	DryRun       bool
}

type tMaxAge struct {
	String   string
	Duration time.Duration
}

func Init(fol, mat, frm, mxa string, rm, yes, skt, dry bool, lg logseal.Logseal) (conf Conf) {
	conf.Now = time.Now()

	abs, err := filepath.Abs(fol)
	lg.IfErrFatal(
		"absolute file path creation failed", logseal.F{
			"error": err,
			"path":  abs,
		},
	)
	if err == nil {
		conf.Folder = abs
	}

	conf.Matcher = mat
	conf.TargetFormat = frm
	conf.MaxAge.String = mxa

	dur, err := str2duration.ParseDuration(mxa)
	lg.IfErrFatal(
		"can not parse max age arg",
		logseal.F{"string": mxa, "error": err},
	)
	if err == nil {
		conf.MaxAge.Duration = dur
	}

	conf.Remove = rm
	conf.Yes = rm
	conf.SkipTruncate = skt
	conf.DryRun = dry
	return
}

func InitTestConf(fol string, remove bool) (conf Conf) {
	lg := logseal.Init()
	conf = Init(
		fol,
		"\\.log$",
		"gz",
		"1s",
		remove,
		false,
		false,
		false,
		lg,
	)
	return
}
