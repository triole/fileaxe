package main

// logTarget returns where the logger should write. Dry runs are exploratory:
// their output is meant to be read on a terminal, so they always log to stdout
// and never into a log file. Every other run honours the requested log file.
func logTarget(logFile string, dryRun bool) string {
	if dryRun {
		return "stdout"
	}
	return logFile
}
