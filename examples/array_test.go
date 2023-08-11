package main

import (
	"fmt"

	"github.com/ivancorrales/knoa"
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
	knoa.Array(knoa.WithStrictMode(true)).Set("[a]", "Jane").JSON()
	// Output:
	// invalid Path  '[a]'. Path doesn't match defined format
}

// Create and array and add/modify entries
func ExampleArrayLoadAndModify() {
	out := knoa.FromArray([]any{"Janet", "Tim"}).Set("[0]", "Jane", "[2]", "Tom").JSON()
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
	k := knoa.Array()
	k.Set("[0][1]", []string{"Tim", "Janet"})
	out := k.JSON()
	fmt.Println(out)
	// Output:
	// [[null,["Tim","Janet"]]]
}

func ExampleArraySetAsteriskAndIndex() {
	initialValue := []any{"red", "blue"}
	k := knoa.FromArray(initialValue)

	k.Set("[2]", "yellow")
	fmt.Println(k.JSON())

	k.Set("[*]", "black")
	fmt.Println(k.JSON())
	// Output:
	// ["red","blue","yellow"]
	// ["black","black","black"]
}

func ExampleArraySetAsteriskAndIndexV2() {
	initialValue := []any{
		Person{
			Firstname: "Jane",
			Age:       20,
		},
		Person{
			Firstname: "Tom",
			Age:       22,
		},
	}
	k := knoa.FromArray(initialValue)

	k.Set("[0].age", 22)
	fmt.Println(k.JSON())

	k.Set("[*].age", 30)
	fmt.Println(k.JSON())
	// Output:
	// [{"age":22,"firstname":"Jane"},{"age":22,"firstname":"Tom"}]
	// [{"age":30,"firstname":"Jane"},{"age":30,"firstname":"Tom"}]
}
