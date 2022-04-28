package bytes

import (
	"github.com/ds248a/lib/strconv"
)

// Copy возврашает копию среза b.
func Copy(b []byte) []byte {
	return []byte(strconv.B2S(b))
}

// Equal проверяе, являются ли преданные срезы a и b эквивалентными.
func Equal(a, b []byte) bool {
	return strconv.B2S(a) == strconv.B2S(b)
}

// Extend расширяет срез b до заданного размера.
func Extend(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}
	return b[:needLen]
}

// Prepend добавляет src в переданный срез dst.
func Prepend(dst []byte, src ...byte) []byte {
	dstLen := len(dst)
	srcLen := len(src)

	dst = Extend(dst, dstLen+srcLen)
	copy(dst[srcLen:], dst[:dstLen])
	copy(dst[:srcLen], src)

	return dst
}

// PrependString добавляет строку в переданный байтовый срез dst.
func PrependString(dst []byte, src string) []byte {
	return Prepend(dst, strconv.S2B(src)...)
}
