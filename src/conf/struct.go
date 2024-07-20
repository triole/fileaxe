package conf

import "time"

type Conf struct {
	Action    string
	Now       time.Time
	Folder    string
	Matcher   string
	MaxAge    time.Duration
	DryRun    bool
	Ls        tLs
	Remove    tRemove
	Rotate    tRotate
	MsgPrefix string
}

type tLs struct {
	Plain bool
}

type tRemove struct {
	Yes bool
}

type tRotate struct {
	CompressionFormat string
	SkipTruncate      bool
}
