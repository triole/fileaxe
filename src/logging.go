package main

import (
	"fileaxe/src/conf"

	"github.com/triole/logseal"
)

// logTarget returns where the logger should write. Dry runs are exploratory:
// their output is meant to be read on a terminal, so they always log to stdout
// and never into a log file. Every other run honours the requested log file.
func logTarget(logFile string, dryRun bool) string {
	if dryRun {
		return "stdout"
	}
	return logFile
}

// logStart reports the run configuration. A plain listing promises file names
// only on stdout, so nothing may be written in front of it.
func logStart(lg logseal.Logseal, cnf conf.Conf) {
	if cnf.Ls.Plain {
		return
	}
	lg.Info("start logaxe", logseal.F{
		"action":  cnf.Action,
		"folder":  cnf.Folder,
		"matcher": cnf.Matcher,
		"minage":  cnf.MinAge,
		"maxage":  cnf.MaxAge,
		"minsize": cnf.MinSize,
		"maxsize": cnf.MaxSize,
	})
}
