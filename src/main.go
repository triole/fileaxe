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
	lg = logseal.Init(
		CLI.LogLevel, logTarget(CLI.LogFile, CLI.DryRun),
		CLI.LogNoColors, CLI.LogJSON,
	)

	conf := conf.Init(CLI, lg)
	la := fileaxe.Init(conf, lg)
	logStart(lg, conf)
	la.Run()
}
