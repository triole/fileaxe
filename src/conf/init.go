package conf

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/triole/logseal"
	str2duration "github.com/xhit/go-str2duration/v2"
)

func Init(cli interface{}, lg logseal.Logseal) (conf Conf) {
	conf.Now = time.Now()

	abs, err := filepath.Abs(getcli(cli, "Folder").(string))
	lg.IfErrFatal(
		"absolute file path creation failed", logseal.F{
			"error": err,
			"path":  abs,
		},
	)
	if err == nil {
		conf.Folder = abs
	}

	conf.Matcher = getcli(cli, "Matcher").(string)
	maxAgeArg := getcli(cli, "MaxAge").(string)

	conf.MaxAge, err = str2duration.ParseDuration(maxAgeArg)
	lg.IfErrFatal(
		"can not parse max age arg",
		logseal.F{"string": maxAgeArg, "error": err},
	)
	conf.DryRun = getcli(cli, "DryRun").(bool)
	if conf.DryRun {
		conf.MsgPrefix = "(dry run) "
	}
	conf.Action = getcli(cli, "SubCommand").(string)
	conf.Ls.Plain = getcli(cli, "Ls.Plain").(bool)
	conf.Remove.Yes = getcli(cli, "Remove.Yes").(bool)
	conf.Rotate.CompressionFormat = getcli(cli, "Rotate.Format").(string)
	conf.Rotate.SkipTruncate = getcli(cli, "Rotate.SkipTruncate").(bool)
	return
}

func getcli(cli interface{}, keypath string) (r interface{}) {
	key := strings.Split(keypath, ".")
	val := reflect.ValueOf(cli)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Name == key[0] {
			r = field.Interface()
			if len(key) > 1 {
				return getcli(r, key[1])
			}
			if fieldType.Type.Name() == "" {
				fmt.Printf("%+v\n", field.Interface())
				r = true
			} else {
				r = field.Interface()
			}
		}
	}
	return
}

func InitTestConf(subcommand, fol string) (conf Conf) {
	conf.Action = subcommand
	conf.Folder = "../testdata/tmp"
	conf.Matcher = ".*"
	conf.MaxAge = 0
	conf.Rotate.CompressionFormat = "gz"
	return
}
