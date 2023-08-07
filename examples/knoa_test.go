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
