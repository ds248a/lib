package rand

import (
	crand "crypto/rand"

	"github.com/ds248a/lib/bpool"
	"github.com/ds248a/lib/bytes"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 62
	charsetIdxBits = 6
	charsetIdxMask = 1<<charsetIdxBits - 1 // 63
)

var randPool = bpool.Pool{}

// Rand заполняет dst случайным набором символов латинского алфавита.
// Необходимое условие: len(dst) > 0
func Rand(dst []byte) {
	if len(dst) == 0 {
		return
	}

	buf := randPool.Get()
	buf.B = bytes.Extend(buf.B, len(dst))

	if _, err := crand.Read(buf.B); err != nil {
		panic(err)
	}

	size := len(dst)
	for i, j := 0, 0; i < size; j++ {
		if idx := int(buf.B[j%size] & charsetIdxMask); idx < len(charset) {
			dst[i] = charset[idx]
			i++
		}
	}

	randPool.Put(buf)
}
