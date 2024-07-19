package main

import (
	"fileaxe/src/conf"
	"fileaxe/src/fileaxe"

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

	la := fileaxe.Init(conf, lg)
	la.Run()
}
