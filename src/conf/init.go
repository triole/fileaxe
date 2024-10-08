package conf

import (
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/triole/logseal"
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
	conf.MinAge, conf.MaxAge = parseDurationRangeArg(
		getcli(cli, "AgeRange").(string), lg,
	)

	conf.SortBy = getcli(cli, "SortBy").(string)
	conf.Order = getcli(cli, "Order").(string)

	conf.DryRun = getcli(cli, "DryRun").(bool)
	if conf.DryRun {
		conf.MsgPrefix = "(dry run) "
	}
	conf.Action = getcli(cli, "SubCommand").(string)
	conf.Ls.Plain = getcli(cli, "List.Plain").(bool)
	conf.Exists.MinNumber, conf.Exists.MaxNumber = parseNumberRangeArg(
		getcli(cli, "Exists.NumberRange").(string), lg,
	)
	conf.Exists.List = getcli(cli, "Exists.List").(bool)
	conf.Remove.Yes = getcli(cli, "Remove.Yes").(bool)
	conf.Compress.CompressionFormat = getcli(cli, "Compress.Format").(string)
	conf.Rotate.CompressionFormat = getcli(cli, "Rotate.Format").(string)
	conf.Rotate.SkipTruncate = getcli(cli, "Rotate.SkipTruncate").(bool)
	conf.Copy.Target = getcli(cli, "Copy.Target").(string)
	conf.Move.Target = getcli(cli, "Move.Target").(string)
	conf.Truncate.Yes = getcli(cli, "Truncate.Yes").(bool)
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
				r = true
			} else {
				r = field.Interface()
			}
		}
	}
	return
}

func InitTestConf(subcommand, folder, matcher string) (conf Conf) {
	conf.Now = time.Now()
	conf.Action = subcommand
	conf.Folder = folder
	conf.Matcher = matcher
	conf.MinAge = 0
	conf.MaxAge = 0
	conf.Rotate.CompressionFormat = "gz"
	conf.Remove.Yes = true
	conf.Truncate.Yes = true
	return
}
