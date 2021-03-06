package str2bytes

import (
	"unsafe"
)

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//go:noinline
func test(b []byte) *int {
	a := 10

	return &a
}