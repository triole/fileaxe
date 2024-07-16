package logaxe

import (
	"logaxe/src/conf"

	"github.com/triole/logseal"
)

type LogAxe struct {
	Conf conf.Conf
	Lg   logseal.Logseal
}

func Init(conf conf.Conf, lg logseal.Logseal) (la LogAxe) {
	la.Conf = conf
	la.Lg = lg
	return
}
