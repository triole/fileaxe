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

	conf := conf.Init(CLI, lg)
	la := fileaxe.Init(conf, lg)
	la.Run()
}
