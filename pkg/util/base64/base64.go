package base64

import "encoding/base64"

//Base64Encode input with base64 encoded.
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

//Base64Decode decode base64-encoded input.
func Base64Decode(input string) []byte {
	decodeBytes, _ := base64.StdEncoding.DecodeString(input)
	return decodeBytes
}
