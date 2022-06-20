package maps

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"testing"
)

func TestFlatten(t *testing.T) {
	json.Print("Nested map: ", nestedMap)
	output := Flatten(nestedMap)
	json.Print("Flatten map", output)
}

func BenchmarkFlatten(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		_ = Flatten(nestedMap)
	}
}

var (
	nestedJson = `{
  "boolean_key": "--- true\n",
  "empty_string_translation": "",
  "key_with_description": "Check it out! This key has a description! (At least in some formats)",
  "key_with_line-break": "This translations contains\na line-break.",
  "nested": {
    "deeply": {
      "key": "Wow, this key is nested even deeper."
    },
    "key": "This key is nested inside a namespace."
  },
  "null_translation": null,
  "pluralized_key": {
    "one": "Only one pluralization found.",
    "other": "Wow, you have %s pluralizations!",
    "zero": "You have no pluralization."
  },
  "sample_collection": [
    "first item",
    "second item",
    "third item"
  ],
  "simple_key": "Just a simple key with a simple message.",
  "unverified_key": "This translation is not yet verified and waits for it. (In some formats we also export this status)"
}`
	nestedMap map[string]interface{}
)

func init() {
	_ = json.Unmarshal([]byte(nestedJson), &nestedMap)
}
