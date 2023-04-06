# Logaxe ![example workflow](https://github.com/triole/example/actions/workflows/build.yaml/badge.svg)

<!--- mdtoc: toc begin -->

1. [Synopsis](#synopsis)
2. [Help](#help)
3. [Disclaimer](#disclaimer)<!--- mdtoc: toc end -->

## Synopsis

Go axing logs. Because I don't like log rotate. Still I little experimental I must admit.

## Help

```go mdox-exec="r -h"

go axing logs

Arguments:
  [<folder>]    folder to process, positional arg required

Flags:
  -h, --help                      Show context-sensitive help.
  -r, --matcher="\\.log$"         regex matcher for file detection
  -m, --max-age="0"               remove compressed log files older than x,
                                  default keeps all, use with duration like:
                                  12h, 4d, 2w, 6m, 2y...
  -l, --log-file="/dev/stdout"    log file
      --log-level="info"          log level
      --log-no-colors             disable output colours, print plain text
      --log-json                  enable json log, instead of text one
  -k, --skip-truncate             skip file truncation, don't empty compressed
                                  log files
  -n, --dry-run                   dry run, just print don't do
  -V, --version-flag              display version
```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
