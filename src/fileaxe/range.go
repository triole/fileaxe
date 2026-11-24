package fileaxe

import "time"

// ageInRange reports whether age fits the age filter. The bounds are exclusive
// and a zero max leaves the range open ended, matching the age range semantics
// of earlier versions.
func ageInRange(age, min, max time.Duration) bool {
	return age > min && (max == 0 || age < max)
}

// sizeInRange reports whether size fits the size filter. Unlike the age filter
// both bounds are inclusive, so -z 1g matches files of 1GiB and above. A zero
// min and max disables the filter, a zero max leaves it open ended.
func sizeInRange(size, min, max int64) bool {
	if min == 0 && max == 0 {
		return true
	}
	if size < min {
		return false
	}
	return max == 0 || size <= max
}

func (fa FileAxe) isInRange(val, min, max int) bool {
	if min == 0 && max == 0 {
		if val == 0 {
			return true
		}
	} else {
		if val >= min && max == 0 {
			return true
		}
		if val >= min && val <= max {
			return true
		}
	}
	return false
}
