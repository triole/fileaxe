package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (la LogAxe) Run() {
	la.Lg.Info(
		"Start fileaxe",
		logseal.F{
			"conf": fmt.Sprintf("%+v", la.Conf),
		},
	)

	if la.Conf.DryRun {
		la.Lg.Info(" --- DRY RUN START ---")
	}

	switch la.Conf.Action {
	case "ls":
		la.list()
	case "rotate":
		la.rotate()
	case "remove":
		la.remove()
	}

	if la.Conf.DryRun {
		la.Lg.Info(" --- DRY RUN END ---")
	}
}
