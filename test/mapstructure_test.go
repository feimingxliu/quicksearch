package test

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"testing"
)

type foo struct {
	A *string `json:"a"`
	B *bar    `json:"b"`
}

type bar struct {
	C *string `json:"c"`
}

func TestMapstructure(t *testing.T) {
	m := map[string]interface{}{
		"a": "a",
		"b": map[string]interface{}{
			"c": "c",
		},
	}
	s := new(foo)
	if err := maps.MapToStruct(m, s); err != nil {
		t.Fatal(err)
	}
	json.Print("struct", s)
}
