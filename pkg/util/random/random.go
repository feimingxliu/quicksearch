package random

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"strings"
)

const (
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32
	digits     = "0123456789"                                           // 10
	characters = letters + digits + symbols                             // 94
)

var buf = make([]byte, 8)

//RandomN generate uint64 random number.
func RandomN(n int) uint64 {
	if n <= 0 {
		log.Fatalln("range of random number <= 0!")
	}
	_, _ = rand.Read(buf)
	rn := binary.BigEndian.Uint64(buf)
	return rn % uint64(n)
}

var stringBuilder = &strings.Builder{}

//RandomString generate random string with length.
func RandomString(length int) string {
	if length <= 0 {
		log.Fatalln("length of random string <= 0!")
	}
	for i := 0; i < length; i++ {
		stringBuilder.WriteByte(RandomChar())
	}
	s := stringBuilder.String()
	stringBuilder.Reset()
	return s
}

//RandomChar generate a random char base on characters above.
func RandomChar() byte {
	return characters[RandomN(len(characters))]
}
