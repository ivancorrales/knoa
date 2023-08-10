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
	// {"age":20,"firstname":"John"}

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
	// {"age":20,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":40,"firstname":"Bob"}]}

	k.Set("age", 23, "siblings[1].age", 39)
	fmt.Println(k.JSON())
	// {"age":23,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":39,"firstname":"Bob"}]}

	k.Set("siblings[*].age", 40)
	fmt.Println(k.JSON())
	// {"age":23,"firstname":"John","siblings":[{"age":40,"firstname":"Tim"},{"age":40,"firstname":"Bob"}]}

	var person Person
	k.To(&person)
	fmt.Println(person)
	// {John 23 [{Tim 40 []} {Bob 40 []}]}
}
