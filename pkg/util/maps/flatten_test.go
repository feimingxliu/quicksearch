package maps

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"testing"
)

func TestFlatten(t *testing.T) {
	inputRaw := `{
	"a" : {
		"b" : [0,1,2],
		"c" : {
				"d": "e"
				}
		},
	"b" : [
			{
				"c" : ["d", "e", "f"]
			},
			{
				"c" : "d",
				"d" : 1
			}
			]
}`
	input := make(map[string]interface{})
	if err := json.Unmarshal([]byte(inputRaw), &input); err != nil {
		t.Fatal(err)
	}
	json.Print("Nested map: ", input)
	output := Flatten(input)
	json.Print("Flatten map", output)
}
