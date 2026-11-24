package fileaxe

import (
	"fmt"

	"github.com/triole/logseal"
)

func (fa FileAxe) Run() {
	if !fa.Conf.Ls.Plain {
		fa.Lg.Debug(
			"start fileaxe",
			logseal.F{"conf": fmt.Sprintf("%+v", fa.Conf)},
		)
	}

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
}
