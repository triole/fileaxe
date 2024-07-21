# Fileaxe ![build](https://github.com/triole/fileaxe/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Help](#help)
- [Disclaimer](#disclaimer)

<!-- /toc -->

## Synopsis

Compress and truncate files that are older than x. Or simply delete them. Can be used as logrotate replacement by just adding a cronjob. Kind of experimental.

## Help

```go mdox-exec="r -h"

find files matching criteria and do something with them

Flags:
  -h, --help                      Show context-sensitive help.
  -f, --folder="/home/ole/rolling/golang/projects/fileaxe/src"
                                  folder to process, default is current
                                  directory
  -m, --matcher="\\..*$"          regex matcher for file detection
  -x, --max-age="0"               max age of files to consider, determined by
                                  last modified date, use with duration like
                                  i.e. 90m, 12h, 4d, 2w
  -s, --sort-by="path"            sort output list by, can be: age, path
  -o, --order="asc"               sort order
      --log-file="/dev/stdout"    log file
      --log-level="info"          log level
      --log-no-colors             disable output colours, print plain text
      --log-json                  enable json log, instead of text one
  -n, --dry-run                   dry run, just print don't do
  -V, --version-flag              display version

Commands:
  ls        list files matching the criteria
  rotate    rotate matching files, compress and truncate after successful
            compression
  remove    remove matching files older than max age

Run "fileaxe <command> --help" for more information on a command.
```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
