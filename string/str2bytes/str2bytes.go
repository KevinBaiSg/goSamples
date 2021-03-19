package str2bytes

import (
	"reflect"
	"unsafe"
)

func str2bytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Cap = sh.Len
	bh.Len = sh.Len
	bh.Data = sh.Data
	return b
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
