package main

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

func ExampleMoney()  {
	pound := money.New(0, "CNY")
	twoPounds := money.New(10000, "CNY")
	// IsZero
	fmt.Println(pound.IsZero())
	// amount
	fmt.Println(twoPounds.Amount())
	// Add
	pound, _ = pound.Add(twoPounds)
	fmt.Println(pound.Display())

	// Output:
	// true
	// 10000
	// 100.00 元
}
