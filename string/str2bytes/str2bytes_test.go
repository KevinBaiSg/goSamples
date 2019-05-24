package str2bytes

import (
	"strings"
	"testing"
)

var s = strings.Repeat("a", 1024)

func testString()  {
	b := []byte(s)
	_ = string(b)
}

func testStr2bytes() {
	b := str2bytes(s)
	_ = bytes2str(b)
}

func BenchmarkTest(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		testString()
	}
}

func BenchmarkTestBlock(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		testStr2bytes()
	}
}
