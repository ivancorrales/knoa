package main

import (
	"fmt"

	"github.com/ivancorrales/knoa"
)

func ExampleArrays() {
	k := knoa.Array().Set("[1]", []string{"red", "blue", "green"}, "[2].firstname", "John")
	fmt.Println(k.String())
	k.Set("[0]", struct {
		FullName  string `structs:"fullName"`
		RoleLevel int    `structs:"roleLevel"`
	}{
		"Senior Developer",
		3,
	})
	fmt.Println(k.String())
	k.Set("[0].enabled", false, "[2].firstname", "Jane")
	fmt.Println(k.String())
	// Output:
	// [null,["red","blue","green"],{"firstname":"John"}]
	// [{"fullName":"Senior Developer","roleLevel":3},["red","blue","green"],{"firstname":"John"}]
	// [{"enabled":false,"fullName":"Senior Developer","roleLevel":3},["red","blue","green"],{"firstname":"Jane"}]
}

type Person struct {
	Firstname string   `structs:"firstname"`
	Age       int      `structs:"age"`
	Siblings  []Person `structs:"siblings,omitempty"`
}

func (p Person) String() string {
	str := fmt.Sprintf("%s -> %d", p.Firstname, p.Age)
	if p.Siblings != nil {
		return fmt.Sprintf("%s siblings[%s]", str, p.Siblings)
	}
	return str
}

func ExampleArrayOfObjects() {
	k := knoa.Load[[]any]([]any{
		Person{
			Firstname: "Jane",
			Age:       20,
		},
	})
	k.Set("[1]", Person{
		Firstname: "Bob",
		Age:       23,
	}, "[2].firstname", "John")
	var output []Person
	err := k.To(&output)
	if err != nil {
		panic(err.Error())
	}
	for i := range output {
		fmt.Println(output[i])
	}
	// Output:
	// Jane -> 20
	// Bob -> 23
	// John -> 0
}

func ExampleMap() {
	k := knoa.Map().Set("firstname", "John", "age", 20)
	fmt.Println(k.String())
	k.Set("siblings", []Person{
		{
			Firstname: "Tim",
			Age:       29,
		}, {
			Firstname: "Bob",
			Age:       40,
		},
	})
	fmt.Println(k.String())
	k.Set("age", 23, "siblings[1].age", 39)
	fmt.Println(k.String())
	var person Person
	k.To(&person)
	fmt.Println(person)

	// Output:
	// {"age":20,"firstname":"John"}
	// {"age":20,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":40,"firstname":"Bob"}]}
	// {"age":23,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":39,"firstname":"Bob"}]}
	// John -> 23 siblings[[Tim -> 29 Bob -> 39]]
}
