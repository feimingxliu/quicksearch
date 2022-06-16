package config

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"testing"
)

func TestConfig(t *testing.T) {
	if err := Init("../../configs/config.yaml"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("", Global)
	}
}
