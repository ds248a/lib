package strconv

import (
	"reflect"
	"unsafe"
)

// B2S преобразует []byte -> string
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B преобразует string -> []byte
func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
