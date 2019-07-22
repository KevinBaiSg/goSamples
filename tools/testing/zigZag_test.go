package testing

import (
	"testing"
	"time"
)

func TestZigZagConversion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cases := []struct {
		s 		string
		numRows int
		want 	string
	}{
		{"", 10, ""},
		{"LEETCODEISHIRING", 0, ""},
		{"LEETCODEISHIRING", 1, "LEETCODEISHIRING"},
		{"LEETCODEISHIRING", 3, "LCIRETOESIIGEDHN"},
		{"LEETCODEISHIRING", 4, "LDREOEIIECIHNTSG"},
	}

	//t.Parallel()

	for _, c := range cases {
		zigZag := convert(c.s, c.numRows)
		if zigZag != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.s, zigZag, c.want)
		}
	}
}

func BenchmarkZigZagConversion(b *testing.B)  {
	//b.ReportAllocs() // 内存统计
	for i := 0; i < b.N; i++ {
		convert("LEETCODEISHIRING", 4)
	}
}

func BenchmarkZigZagConversionResetTimer(b *testing.B)  {

	time.Sleep(100 * time.Millisecond)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convert("LEETCODEISHIRING", 4)
	}
}

func BenchmarkZigZagConversionRunParallel(b *testing.B)  {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			convert("LEETCODEISHIRING", 4)
		}
	})
}
//
//func ExampleZigZagConversion() {
//	convert("LEETCODEISHIRING", 4)
//	// Output: LDREOEIIECIHNTSG
//}