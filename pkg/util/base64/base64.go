package base64

import "encoding/base64"

//Encode input with base64 encoded.
func Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

//Decode decode base64-encoded input.
func Decode(input string) []byte {
	decodeBytes, _ := base64.StdEncoding.DecodeString(input)
	return decodeBytes
}
