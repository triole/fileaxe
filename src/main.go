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
	la.Lg.Info("start logaxe", logseal.F{
		"action":  conf.Action,
		"folder":  conf.Folder,
		"matcher": conf.Matcher,
		"minage":  conf.MinAge,
		"maxage":  conf.MaxAge,
		"minsize": conf.MinSize,
		"maxsize": conf.MaxSize,
	})
	la.Run()
}
