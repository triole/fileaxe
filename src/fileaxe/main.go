package fileaxe

import (
	"time"

	"github.com/triole/logseal"
)

func (fa FileAxe) Run() {
	start := time.Now()
	fa.logStart()

	fileList := fa.Find(Criteria{
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
		fa.exists(fileList)
	case "compress":
		fa.compress(fileList)
	case "rotate":
		fa.rotate(fileList)
	case "copy":
		fa.copy(fileList)
	case "move":
		fa.move(fileList)
	case "truncate":
		fa.truncate(fileList)
	case "remove":
		fa.remove(fileList)
	}
	fa.logDone(fileList, time.Since(start))
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
		"min_size":  fa.Conf.MinSize,
		"max_size":  fa.Conf.MaxSize,
	})
}

// logDone reports the outcome of a run: how long it took and what the matched
// set looked like. Like the start message it stays silent during a plain
// listing, whose output is reserved for file names. Oldest and newest age are
// omitted for an empty match, where they would suggest a file was found.
func (fa FileAxe) logDone(fileList FileInfos, elapsed time.Duration) {
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
	fa.Lg.Info("done fileaxe", fields)
}
