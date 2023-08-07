package main

import (
	"fmt"

	"github.com/ivancorrales/knoa"
	"github.com/ivancorrales/knoa/mapifier"
)

// Create and array and add an entry
func ExampleArrayFromScratch() {
	out := knoa.Array().Set("[0]", "Jane").JSON()
	fmt.Println(out)
	// Output:
	// ["Jane"]
}

// Panic when one or more of the passed paths is not valid. by default It's false
func ExampleArrayStrict() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	knoa.Array(mapifier.WithStrictMode(true)).Set("[a]", "Jane").JSON()
	// Output:
	// invalid Path  '[a]'. Path doesn't match defined format
}

// Create and array and add/modify entries
func ExampleArrayLoadAndModify() {
	out := knoa.Load([]any{"Janet", "Tim"}).Set("[0]", "Jane", "[2]", "Tom").JSON()
	fmt.Println(out)
	// Output:
	// ["Jane","Tim","Tom"]
}

// Set some invalid indexes
func ExampleArraySetInvalidIndexes() {
	out := knoa.Array().Set("person-firstname", "Jane", "[2]", true, "lastname", "Doe").JSON()
	fmt.Println(out)
	// Output:
	//[null,null,true]
}

// Set an array as the value of an attribute
func ExampleArraySetSubArrays() {
	out := knoa.Array().Set("[0]", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// [["Tim","Janet"]]
}

// Set values for two-deep level of arrays attributes
func ExampleArraySetSubArraysV2() {
	out := knoa.Array().Set("[0][1]", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// [[null,["Tim","Janet"]]]
}
