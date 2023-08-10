package main

import (
	"fmt"

	"github.com/ivancorrales/knoa"
)

type Person struct {
	Firstname string   `structs:"firstname"`
	Age       int      `structs:"age"`
	Siblings  []Person `structs:"siblings,omitempty"`
}

func main() {
	k := knoa.Map().Set("firstname", "John", "age", 20)
	fmt.Println(k.JSON())
	// fmt.Println(k.YAML())
	// fmt.Println(k.Out())

	k.Set("siblings", []Person{
		{
			Firstname: "Tim",
			Age:       29,
		},
		{
			Firstname: "Bob",
			Age:       40,
		},
	})
	fmt.Println(k.JSON())
	// fmt.Println(k.YAML())
	// fmt.Println(k.Out())

	k.Set("age", 23, "siblings[1].age", 39)
	fmt.Println(k.JSON())
	// fmt.Println(k.YAML())
	// fmt.Println(k.Out())

	k.Set("siblings[*].age", 40)
	fmt.Println(k.JSON())
	// fmt.Println(k.YAML())
	// fmt.Println(k.Out())

	var person Person
	k.To(&person)
	fmt.Println(person)
}
