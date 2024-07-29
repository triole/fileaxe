package fileaxe

import (
	"fileaxe/src/conf"
	"testing"

	"github.com/triole/logseal"
)

func TestIsInRange(t *testing.T) {
	validateIsInRange(3, 1, 5, true, t)
	validateIsInRange(0, 0, 0, true, t)
	validateIsInRange(1, 0, 0, false, t)
	validateIsInRange(1, 1, 0, true, t)
	validateIsInRange(20, 1, 0, true, t)
	validateIsInRange(2, 1, 0, true, t)
	validateIsInRange(2, 1, 1, false, t)
}

func validateIsInRange(val, min, max int, exp bool, t *testing.T) {
	conf := conf.InitTestConf("ex", "../testdata/tmp", "\\..*$")
	lg := logseal.Init()
	fa := Init(conf, lg)
	res := fa.isInRange(val, min, max)
	if res != exp {
		t.Errorf(
			"test isInRange failed: val, min, max: %d, %d, %d -- exp: %v, res: %v",
			val, min, max, exp, res,
		)
	}
}
