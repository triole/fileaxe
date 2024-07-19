package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (la LogAxe) Run() {
	la.Lg.Info(
		"Start logaxe",
		logseal.F{"conf": fmt.Sprintf("%+v", la.Conf)},
	)

	if la.Conf.DryRun {
		la.Lg.Info(" --- DRY RUN START ---")
	}

	if la.Conf.Remove {
		la.remove()
	} else {
		la.compress()
	}

	if la.Conf.DryRun {
		la.Lg.Info(" --- DRY RUN END ---")
	}
}
