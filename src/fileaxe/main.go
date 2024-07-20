package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (fa FileAxe) Run() {
	if !fa.Conf.Ls.Plain {
		fa.Lg.Info(
			"start fileaxe",
			logseal.F{
				"conf": fmt.Sprintf("%+v", fa.Conf),
			},
		)

		if fa.Conf.DryRun {
			fa.Lg.Info(" --- DRY RUN START ---")
		}
	}

	switch fa.Conf.Action {
	case "ls":
		fa.list()
	case "rotate":
		fa.rotate()
	case "remove":
		fa.remove()
	}

	if !fa.Conf.Ls.Plain {
		if fa.Conf.DryRun {
			fa.Lg.Info(" --- DRY RUN END ---")
		}
	}
}
