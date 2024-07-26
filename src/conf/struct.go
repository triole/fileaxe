package conf

import "time"

type Conf struct {
	Action    string
	Now       time.Time
	Folder    string
	Matcher   string
	MinAge    time.Duration
	MaxAge    time.Duration
	SortBy    string
	Order     string
	DryRun    bool
	Ls        tLs
	Rotate    tRotate
	Move      tMove
	Remove    tRemove
	MsgPrefix string
}

type tLs struct {
	Plain bool
}

type tMove struct {
	Target string
}

type tRemove struct {
	Yes bool
}

type tRotate struct {
	CompressionFormat string
	SkipTruncate      bool
}
