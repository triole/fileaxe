package conf

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/triole/logseal"
)

// sizeUnits maps a unit suffix to its multiplier in bytes. Sizes are binary
// units (1k = 1024 bytes), matching how file sizes are usually reported. A
// trailing "b" or "ib" is accepted, so 1k, 1kb and 1kib all mean 1024 bytes.
var sizeUnits = map[string]int64{
	"": 1, "b": 1,
	"k": 1 << 10, "kb": 1 << 10, "kib": 1 << 10,
	"m": 1 << 20, "mb": 1 << 20, "mib": 1 << 20,
	"g": 1 << 30, "gb": 1 << 30, "gib": 1 << 30,
	"t": 1 << 40, "tb": 1 << 40, "tib": 1 << 40,
	"p": 1 << 50, "pb": 1 << 50, "pib": 1 << 50,
}

// parseSize converts a size string such as "1024", "10k", "1.5g" or "2 GiB"
// into a byte count. Suffixes are case insensitive, negative values are
// rejected.
func parseSize(s string) (int64, error) {
	str := strings.TrimSpace(s)
	if str == "" {
		return 0, fmt.Errorf("empty size")
	}

	numPart, unitPart := str, ""
	if idx := strings.IndexFunc(str, func(r rune) bool {
		return (r < '0' || r > '9') && r != '.'
	}); idx >= 0 {
		numPart, unitPart = str[:idx], str[idx:]
	}

	mult, found := sizeUnits[strings.ToLower(strings.TrimSpace(unitPart))]
	if !found {
		return 0, fmt.Errorf("unknown size unit %q", strings.TrimSpace(unitPart))
	}

	val, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}
	if val < 0 || math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("invalid size %q", s)
	}

	bytes := val * float64(mult)
	if bytes > math.MaxInt64 {
		return 0, fmt.Errorf("size %q is too large", s)
	}
	return int64(math.Round(bytes)), nil
}

// parseSizeRangeArg parses a size range of one or two comma separated sizes.
// Both bounds are inclusive, a missing or zero max leaves the range open ended
// and 0,0 disables the filter.
func parseSizeRangeArg(s string, lg logseal.Logseal) (min, max int64) {
	var err error
	arr := strings.Split(s, ",")
	min, err = parseSize(arr[0])
	lg.IfErrFatal(
		"can not parse size range arg",
		logseal.F{"string": arr[0], "error": err},
	)
	if len(arr) > 1 {
		max, err = parseSize(arr[1])
		lg.IfErrFatal(
			"can not parse size range arg",
			logseal.F{"string": arr[1], "error": err},
		)
	}
	if max > 0 && max < min {
		lg.IfErrFatal(
			"can not parse size range arg",
			logseal.F{
				"string": s,
				"error":  fmt.Errorf("max size is smaller than min size"),
			},
		)
	}
	return
}
