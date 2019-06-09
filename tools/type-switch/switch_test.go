package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var s = strings.Repeat("a", 1024)

func testUnmarshalInterface()  {
	type Message struct {
		Name string
		Time interface{}
	}

	b := []byte(`{"Name":"Wednesday","Time":6}`)

	var f Message
	json.Unmarshal(b, &f)

	switch v := f.Time.(type) {
	case float64:
		fmt.Printf("Name: %v; Time: %v\n", f.Name, v)
		fmt.Printf("Time type %T!\n", v)
	case string:
		fmt.Printf("Name: %v; Time: %v\n", f.Name, v)
		fmt.Printf("Time type %T!\n", v)
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

func testUnmarshal() {
	type Message struct {
		Name string
		Time float64
	}

	b := []byte(`{"Name":"Wednesday","Time":6}`)

	var f Message
	json.Unmarshal(b, &f)

	fmt.Printf("Name: %v; Time: %v\n", f.Name, f.Time)
	fmt.Printf("Time type %T!\n", f.Time)
}

func BenchmarkTestUnmarshalInterface(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		testUnmarshalInterface()
	}
}

func BenchmarkTestUnmarshal(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		testUnmarshal()
	}
}