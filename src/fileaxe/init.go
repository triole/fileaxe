package fileaxe

import (
	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

type FileAxe struct {
	Conf conf.Conf
	Lg   logseal.Logseal
}

func Init(conf conf.Conf, lg logseal.Logseal) (fa FileAxe) {
	fa.Conf = conf
	fa.Lg = lg
	return
}
