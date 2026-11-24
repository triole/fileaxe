package fileaxe

import (
	"errors"
	"time"

	"github.com/triole/logseal"
)

// Run executes the configured action. It returns an error when any file could
// not be processed, so main can signal the failure through the exit code; the
// completion event carries the number of failures.
func (fa FileAxe) Run() error {
	start := time.Now()
	fa.logStart()

	fileList, errs := fa.Find(Criteria{
		Folder:  fa.Conf.Folder,
		Matcher: fa.Conf.Matcher,
		MinAge:  fa.Conf.MinAge,
		MaxAge:  fa.Conf.MaxAge,
		MinSize: fa.Conf.MinSize,
		MaxSize: fa.Conf.MaxSize,
		SortBy:  fa.Conf.SortBy,
		Order:   fa.Conf.Order,
		RefTime: fa.Conf.Now,
	})
	switch fa.Conf.Action {
	case "list":
		fa.list(fileList)
	case "exists":
		errs = errors.Join(errs, fa.exists(fileList))
	case "compress":
		errs = errors.Join(errs, fa.compress(fileList))
	case "rotate":
		errs = errors.Join(errs, fa.rotate(fileList))
	case "copy":
		errs = errors.Join(errs, fa.copy(fileList))
	case "move":
		errs = errors.Join(errs, fa.move(fileList))
	case "truncate":
		errs = errors.Join(errs, fa.truncate(fileList))
	case "remove":
		errs = errors.Join(errs, fa.remove(fileList))
	}
	fa.logDone(fileList, time.Since(start), countErrors(errs))
	return errs
}

// logStart reports the run configuration once, before any file is touched. A
// plain listing promises file names only on stdout, so it stays silent there.
// The field names match the events emitted by Find, so a log query for min_age
// finds both the start line and the per-file decisions.
func (fa FileAxe) logStart() {
	if fa.Conf.Ls.Plain {
		return
	}
	fa.Lg.Info("start fileaxe", logseal.F{
		"action":    fa.Conf.Action,
		"folder":    fa.Conf.Folder,
		"rxmatcher": fa.Conf.Matcher,
		"symlinks":  fa.Conf.FollowSymlinks,
		"min_age":   humanReadableDuration(fa.Conf.MinAge),
		"max_age":   humanReadableDuration(fa.Conf.MaxAge),
		"min_size":  humanReadableFileSize(float64(fa.Conf.MinSize)),
		"max_size":  humanReadableFileSize(float64(fa.Conf.MaxSize)),
	})
}

// logDone reports the outcome of a run: how long it took and what the matched
// set looked like. Like the start message it stays silent during a plain
// listing, whose output is reserved for file names. Oldest and newest age are
// omitted for an empty match, where they would suggest a file was found, and
// the error count is omitted for a clean run.
func (fa FileAxe) logDone(fileList FileInfos, elapsed time.Duration, errCount int) {
	if fa.Conf.Ls.Plain {
		return
	}
	var totalSize int64
	var oldest, newest time.Duration
	for i, fil := range fileList {
		totalSize += fil.Size
		if i == 0 || fil.Age > oldest {
			oldest = fil.Age
		}
		if i == 0 || fil.Age < newest {
			newest = fil.Age
		}
	}
	fields := logseal.F{
		"action":   fa.Conf.Action,
		"no":       len(fileList),
		"size":     humanReadableFileSize(float64(totalSize)),
		"duration": humanReadableDuration(elapsed),
	}
	if len(fileList) > 0 {
		fields["oldest_age"] = humanReadableDuration(oldest)
		fields["newest_age"] = humanReadableDuration(newest)
	}
	if errCount > 0 {
		fields["errors"] = errCount
	}
	fa.Lg.Info("done fileaxe", fields)
}

// countErrors returns how many errors were joined into err. The completion
// event reports a count instead of repeating every error the actions logged.
func countErrors(err error) (n int) {
	if err == nil {
		return 0
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range joined.Unwrap() {
			n += countErrors(e)
		}
		return n
	}
	return 1
}
