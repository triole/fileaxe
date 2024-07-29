package conf

import (
	"testing"
	"time"

	"github.com/triole/logseal"
)

func TestParseDurationRangeArg(t *testing.T) {
	lg := logseal.Init()
	validateParseDurationRangeArg("10s", time.Second*10, time.Second*0, lg, t)
	validateParseDurationRangeArg("10s,0", time.Second*10, time.Second*0, lg, t)
	validateParseDurationRangeArg("120s,1h", time.Second*120, time.Second*3600, lg, t)
	validateParseDurationRangeArg("0,30m", time.Second*0, time.Second*1800, lg, t)
	validateParseDurationRangeArg("0", time.Second*0, time.Second*0, lg, t)
}

func TestParseNumberRangeArg(t *testing.T) {
	lg := logseal.Init()
	validateParseNumberRangeArg("1", 1, 1, lg, t)
	validateParseNumberRangeArg("2,0", 2, 0, lg, t)
	validateParseNumberRangeArg("11,0", 11, 0, lg, t)
	validateParseNumberRangeArg("3,30", 3, 30, lg, t)
	validateParseNumberRangeArg("4,40", 4, 40, lg, t)
	validateParseNumberRangeArg("0,50", 0, 50, lg, t)
}

func validateParseDurationRangeArg(s string, expMin, expMax time.Duration, lg logseal.Logseal, t *testing.T) {
	min, max := parseDurationRangeArg(s, lg)
	if min != expMin {
		t.Errorf(
			"failed test parseDurationRangeArg, min != expMin: %s != %s",
			min, expMin,
		)
	}
	if max != expMax {
		t.Errorf(
			"failed test parseDurationRangeArg, max != expMax: %s != %s",
			max, expMax,
		)
	}
}

func validateParseNumberRangeArg(s string, expMin, expMax int, lg logseal.Logseal, t *testing.T) {
	min, max := parseNumberRangeArg(s, lg)
	if min != expMin {
		t.Errorf(
			"failed test parseDurationRangeArg, min != expMin: %d != %d",
			min, expMin,
		)
	}
	if max != expMax {
		t.Errorf(
			"failed test parseDurationRangeArg, max != expMax: %d != %d",
			max, expMax,
		)
	}
}
