// Package util is for utilities
package util

import "bytes"

func ClearZeros(b []byte) []byte {
	i := bytes.IndexByte(b, 0)
	if i < 0 {
		i = len(b)
	}

	return b[:i]
}
