package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "fileaxe"
	appDescription = "find files matching criteria and do something with them"
	appMainversion = "0.4"
)

var CLI struct {
	SubCommand  string `kong:"-"`
	Folder      string `help:"folder to process, default is current directory" short:"f" default:"${curdir}"`
	Matcher     string `help:"regex matcher for file detection" short:"m" default:"\\..*$"`
	MaxAge      string `help:"max age of files to consider, determined by last modified date, use with duration like i.e. 90m, 12h, 4d, 2w" short:"x" default:"0"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level" default:"info" enum:"trace,debug,info,error,fatal"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	LogJSON     bool   `help:"enable json log, instead of text one"`
	DryRun      bool   `help:"dry run, just print don't do" short:"n"`
	VersionFlag bool   `help:"display version" short:"V"`

	Ls struct {
		Plain bool `help:"print plain list, file names only" short:"p"`
	} `cmd:"" help:"list files matching the criteria"`

	Rotate struct {
		Format       string `help:"compression format, if files are not removed" short:"g" default:"gz" enum:"snappy,gz,xz"`
		SkipTruncate bool   `help:"skip file truncation, don't empty compressed log files" short:"k"`
	} `cmd:"" help:"rotate matching files, compress and truncate after successful compression"`

	Remove struct {
		Yes bool `help:"assume yes on remove affirmation query"`
	} `cmd:"" help:"remove matching files older than max age"`
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
		},
	)
	switch ctx.Command() {
	case "ls":
		CLI.SubCommand = "ls"
	case "rotate":
		CLI.SubCommand = "rotate"
	case "remove":
		CLI.SubCommand = "remove"
	default:
		panic(ctx.Command())
	}
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
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
