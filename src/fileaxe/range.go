package fileaxe

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
