package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "logaxe"
	appDescription = "go axing logs"
	appMainversion = "0.1"
)

var CLI struct {
	Folder       string `help:"folder to process, positional arg required" arg:"" optional:""`
	Matcher      string `help:"regex matcher for file detection" short:"r" default:"\.log$"`
	MaxAge       string `help:"remove compressed log files older than x, default keeps all, use with duration like i.e. 90m, 12h, 4d, 2w" short:"m" default:"0"`
	LogFile      string `help:"log file" short:"l" default:"/dev/stdout"`
	LogLevel     string `help:"log level" default:"info" enum:"trace,debug,info,error,fatal"`
	LogNoColors  bool   `help:"disable output colours, print plain text"`
	LogJSON      bool   `help:"enable json log, instead of text one"`
	SkipTruncate bool   `help:"skip file truncation, don't empty compressed log files" short:"k"`
	DryRun       bool   `help:"dry run, just print don't do" short:"n"`
	VersionFlag  bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"curdir": curdir,
			"config": path.Join(getBindir(), appName+".toml"),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	if CLI.Folder == "" {
		fmt.Printf("%s\n", "[ERROR] positional arg required, please pass folder name")
		os.Exit(1)
	}
	// ctx.FatalIfErrorf(err)
}

type tPrinter []tPrinterEl
type tPrinterEl struct {
	Key string
	Val string
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "version: "+appMainversion+".", -1)
	fmt.Printf("\n%s\n%s\n\n", appName, appDescription)
	arr := strings.Split(s, "\n")
	var pr tPrinter
	var maxlen int
	for _, line := range arr {
		if strings.Contains(line, ":") {
			l := strings.Split(line, ":")
			if len(l[0]) > maxlen {
				maxlen = len(l[0])
			}
			pr = append(pr, tPrinterEl{l[0], strings.Join(l[1:], ":")})
		}
	}
	for _, el := range pr {
		fmt.Printf("%"+strconv.Itoa(maxlen)+"s\t%s\n", el.Key, el.Val)
	}
	fmt.Printf("\n")
}

func getBindir() (s string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	s = filepath.Dir(ex)
	return
}
