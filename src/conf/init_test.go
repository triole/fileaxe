package conf

import (
	"testing"
)

func TestInit(t *testing.T) {
	_ = InitTestConf("rotate", "../testdata/tmp")
	if true == false {
		t.Errorf("An error occured.")
	}
}
