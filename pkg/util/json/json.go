package json

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
)

var (
	j = jsoniter.ConfigCompatibleWithStandardLibrary
)

//Print print v in json format, prefix attached at head.
func Print(prefix string, v interface{}) {
	b, _ := j.MarshalIndent(v, "", " ")
	fmt.Printf("%s:\n%s\n", prefix, b)
}

//Marshal same as std json.Marshal, but more efficient.
func Marshal(v interface{}) ([]byte, error) {
	b, err := j.Marshal(v)
	return b, err
}

//NewEncoder same as std json.NewEncoder.
func NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return j.NewEncoder(writer)
}

//Unmarshal same as std json.Unmarshal, but more efficient.
func Unmarshal(data []byte, v interface{}) error {
	return j.Unmarshal(data, v)
}

//NewDecoder same as std json.NewDecoder.
func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return j.NewDecoder(reader)
}
