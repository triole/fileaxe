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
	AgeRange    string `help:"age range of files to consider, string of one or two comma separated values (min age and max age), supports durations like 90m, 12h, 4d, 2w; default behaviour is that all files in a folder will be considered, usage: -r 2h, -r 30m,2h" short:"r" default:"0,0"`
	SortBy      string `help:"sort output list by, can be: age, path" short:"s" enum:"age,path" default:"age"`
	Order       string `help:"sort order" short:"o" enum:"asc,desc" default:"desc"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level" default:"info" enum:"trace,debug,info,error,fatal"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	LogJSON     bool   `help:"enable json log, instead of text one"`
	DryRun      bool   `help:"dry run, just print don't do" short:"n"`
	VersionFlag bool   `help:"display version" short:"V"`

	Ls struct {
		Plain bool `help:"print plain list, file names only" short:"p"`
	} `cmd:"" help:"list files matching the criteria"`

	Ex struct {
		NumberRange string `help:"number of files to be considered a valid match, check is successful if the number of matched files is in the expected range, arg is a string of one or two comma separated values (min and max), e.g. '1' requires exactly one match, '1,5' represents the range between 1 and 5, '1,0' is default meaning any number of matches higher than one will do" short:"b" default:"1,0"`
	} `cmd:"" help:"check if file(s) exists, return non-zero exitcode if not"`

	Rt struct {
		Format       string `help:"compression format, if files are not removed" short:"g" default:"gz" enum:"snappy,gz,xz"`
		SkipTruncate bool   `help:"skip file truncation, don't empty compressed log files" short:"k"`
	} `cmd:"" help:"rotate matching files, compress and truncate after successful compression"`

	Cp struct {
		Target string `help:"target folder to which the files are copied" short:"t" required:""`
	} `cmd:"" help:"copy matching files, requires target folder definition"`

	Mv struct {
		Target string `help:"target folder to which the files are moved" short:"t" required:""`
	} `cmd:"" help:"move matching files, requires target folder definition"`

	Tn struct {
		Yes bool `help:"assume yes on truncate affirmation query"`
	} `cmd:"" help:"truncate matching files"`

	Rm struct {
		Yes bool `help:"assume yes on remove affirmation query"`
	} `cmd:"" help:"remove matching files"`
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
	CLI.SubCommand = ctx.Command()
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
