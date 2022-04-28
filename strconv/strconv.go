package strconv

import (
	"unsafe"
)

// B2S преобразует []byte -> string
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B преобразует string -> []byte
func S2B(s string) (b []byte) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
