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
	Exists    tExists
	Rotate    tRotate
	Copy      tCopyMove
	Move      tCopyMove
	Truncate  tRemove
	Remove    tRemove
	MsgPrefix string
}

type tLs struct {
	Plain bool
}

type tExists struct {
	MinNumber int
	MaxNumber int
	List      bool
}

type tCopyMove struct {
	Target string
}

type tRemove struct {
	Yes bool
}

type tRotate struct {
	CompressionFormat string
	SkipTruncate      bool
}
