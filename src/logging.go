package main

import (
	"os"

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

// isStdoutTarget reports whether a log target means standard output, matching
// how logseal interprets the value it is given.
func isStdoutTarget(target string) bool {
	return target == "" || target == "stdout" || target == "/dev/stdout"
}

// routeLogsForPlainListing keeps stdout clean for a plain listing, which
// promises file names only. When the log stream would go to stdout, it moves to
// stderr; an explicit log file keeps receiving the logs.
func routeLogsForPlainListing(lg logseal.Logseal, target string, plain bool) {
	if !plain || !isStdoutTarget(target) {
		return
	}
	lg.Logrus.SetOutput(os.Stderr)
}
