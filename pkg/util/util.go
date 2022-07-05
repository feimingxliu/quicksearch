package util

import (
	"io"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"strings"
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

// CopyDir copies the content of src to dst. src should be a full path.
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// copy to this path
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			_ = os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}

		// handle irregular files
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		// copy contents of regular file efficiently

		// open input
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// make it the same
		_ = fh.Chmod(info.Mode())

		// copy content
		_, err = io.Copy(fh, in)
		return err
	})
}

func FileExists(path string) (bool, error) {
	if _, err := os.Open(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
