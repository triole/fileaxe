package main

import (
	"logaxe/src/conf"
	"logaxe/src/logaxe"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	conf := conf.Init(
		CLI.Folder,
		CLI.Matcher,
		CLI.Format,
		CLI.MaxAge,
		CLI.Remove,
		CLI.Yes,
		CLI.SkipTruncate,
		CLI.DryRun,
		lg,
	)

	la := logaxe.Init(conf, lg)
	la.Run()
}
