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

// Ignore those attributes that don't match the provided format
func ExampleArrayWithAttributeNameFormat() {
	out := knoa.Array(mapifier.WithAttributeNameFormat("person-(.*)")).Set("person-firstname", "Jane", "lastname", "Doe").JSON()
	fmt.Println(out)
	// Output:
	// {"person-firstname":"Jane"}
}

// Set an array as the value of an attribute
func ExampleSArraySetSubarrays() {
	out := knoa.Array().Set("[0]", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// [["Tim","Janet"]]
}

// Set values for two-deep level of arrays attributes
func ExampleSArraySetSubarraysV2() {
	out := knoa.Array().Set("[0][1]", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// [[],[["Tim","Janet"]]]
}
