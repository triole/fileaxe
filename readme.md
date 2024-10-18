# Fileaxe ![build](https://github.com/triole/fileaxe/actions/workflows/build.yaml/badge.svg) ![test](https://github.com/triole/fileaxe/actions/workflows/test.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Help](#help)
- [Disclaimer](#disclaimer)

<!-- /toc -->

## Synopsis

A file finder that takes the last modified date into consideration. It can list files in a certain time span or use them for different kinds of processing actions.

## Help

```go mdox-exec="r -h"

find files matching criteria and do something with them

Flags:
  -h, --help                      Show context-sensitive help.
  -f, --folder="/home/ole/rolling/golang/projects/fileaxe/src"
                                  folder to process, default is current
                                  directory
  -m, --matcher=".*"              regex matcher for file detection, e.g. '\..*$'
                                  or '\.yaml$
  -r, --age-range="0,0"           age range of files to consider, string of one
                                  or two comma separated values (min age and
                                  max age), supports durations like 90m, 12h,
                                  4d, 2w; default behaviour is that all files
                                  in a folder will be considered, usage: -r 2h,
                                  -r 30m,2h
  -s, --sort-by="age"             sort output list by, can be: age, path
  -o, --order="desc"              sort order
      --log-file="/dev/stdout"    log file
      --log-level="info"          log level
      --log-no-colors             disable output colours, print plain text
      --log-json                  enable json log, instead of text one
  -n, --dry-run                   dry run, just print don't do

Commands:
  list        list files matching the criteria
  exists      check if file(s) exists, return non-zero exitcode if not
  compress    compress matching files
  rotate      rotate matching files, compress and truncate after successful
              compression
  copy        copy matching files, requires target folder definition
  move        move matching files, requires target folder definition
  truncate    truncate matching files
  remove      remove matching files
  version     display version

Run "fileaxe <command> --help" for more information on a command.
```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
