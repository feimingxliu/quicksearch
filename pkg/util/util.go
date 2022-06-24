package util

import (
	"math/big"
	"time"
)

func ExecTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func BytesModInt(b []byte, i int) int64 {
	if i == 0 {
		return 0
	}
	bi := big.Int{}
	return bi.Mod((&big.Int{}).SetBytes(b), big.NewInt(int64(i))).Int64()
}
