package util

import (
	"encoding/hex"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"testing"
)

//go test -v github.com/feimingxliu/quicksearch/pkg/util -run 'BytesModInt' -count 1
func TestBytesModInt(t *testing.T) {
	for i := 1; i <= 10; i++ {
		b, err := hex.DecodeString(uuid.GetUUID())
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%x mod %d: %d\n", b, 10, BytesModInt(b, 10))
	}
}

//go test -v github.com/feimingxliu/quicksearch/pkg/util -bench 'BytesModInt' -benchmem
func BenchmarkBytesModInt(b *testing.B) {
	m := make(map[int64]uint)
	for i := 0; i < b.N; i++ {
		bs, err := hex.DecodeString(uuid.GetUUID())
		if err != nil {
			b.Fatal(err)
		}
		m[BytesModInt(bs, 10)]++
	}
	for k, v := range m {
		fmt.Printf("%d: %d times.\n", k, v)
	}
}
