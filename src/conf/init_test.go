package conf

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	conf := InitTestConf("../testdata/tmp", false)
	fmt.Printf("%+v\n", conf)
	if true == false {
		t.Errorf("An error occured.")
	}
}
