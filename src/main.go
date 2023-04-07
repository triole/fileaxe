package main

import (
	"logaxe/src/logaxe"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()

	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	la := logaxe.InitLogAxe(
		CLI.Folder, CLI.Matcher, CLI.MaxAge, CLI.SkipTruncate, CLI.DryRun, lg,
	)

	la.Run()
}
