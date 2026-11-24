package main

import (
	"os"

	"fileaxe/src/conf"
	"fileaxe/src/fileaxe"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	target := logTarget(CLI.LogFile, CLI.DryRun)
	lg = logseal.Init(CLI.LogLevel, target, CLI.LogNoColors, CLI.LogJSON)
	// CLI.List.Plain is what conf.Init copies into Conf.Ls.Plain. Routing early
	// keeps even argument errors off the plain listing's stdout.
	routeLogsForPlainListing(lg, target, CLI.List.Plain)

	conf := conf.Init(CLI, lg)
	la := fileaxe.Init(conf, lg)
	if err := la.Run(); err != nil {
		os.Exit(1)
	}
}
