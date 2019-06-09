package main

import (
	"encoding/json"
	"fmt"
)

func main() {
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

	b = []byte(`{"Name":"Wednesday","Time":"6"}`)
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
