package main

import (
	"fmt"
	"strings"

	"github.com/ivancorrales/mapify"
	"github.com/ivancorrales/mapify/mapifier"
)

// Basic showcase
func ExampleFromScratch() {
	out := mapify.Map().Set("firstname", "Jane").JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane"}
}

// Panic when the attribute is not valid
func ExampleFromScratchWithStrictModeEnabled() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	mapify.Map(mapifier.WithStrictMode(true)).Set("firstname.$", "Jane").JSON()
	// Output:
	// invalid Path  'firstname.$'. Path doesn't match defined format
}

// Ignore those attributes that don't match the provided format
func ExampleFromScratchWithAttributeNameFormat() {
	out := mapify.Map(mapifier.WithAttributeNameFormat("person-(.*)")).Set("person-firstname", "Jane", "lastname", "Doe").JSON()
	fmt.Println(out)
	// Output:
	// {"person-firstname":"Jane"}
}

// Set arrays attributes
func ExampleSetArrayChildren() {
	out := mapify.Map().Set("firstname", "Jane", "siblings", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","siblings":["Tim","Janet"]}
}

// Set arrays attributes
func ExampleSetArrayChildrenV2() {
	out := mapify.Map().Set("firstname", "Jane", "languages.native", []string{"English", "Irish"}, "languages.learning", []string{"Italian"}).JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","languages":{"learning":["Italian"],"native":["English","Irish"]}}
}

// Set several times
func ExampleMultipleSets() {
	m := mapify.Map().Set("firstname", "Jane")
	m = m.Set("lastname", "Doe")
	m = m.Set("firstname", "Tim")
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Tim","lastname":"Doe"}
}

// Set complex structures
func ExampleSetComplexStructures() {
	out := mapify.Map().Set("firstname", "Jane", "partner", struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		32, "Tim",
	}).JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","partner":{"age":32,"firstname":"Tim"}}
}

// Set complex structures
func ExampleSetComplexStructuresAndOverrides() {
	m := mapify.Map().Set("firstname", "Jane", "partner", struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		32, "Tim",
	})
	m = m.Set("partner.age", 40)

	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","partner":{"age":40,"firstname":"Tim"}}
}

func ExampleWithFuncPrefix() {
	m := mapify.Map().Set("firstname", "Jane", "partner", struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		32, "Tim",
	})
	m = m.With(mapifier.WithFuncPrefix(strings.ToUpper))("gender", "female")

	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"GENDER":"female","firstname":"Jane","partner":{"age":32,"firstname":"Tim"}}
}

func ExampleWithPrefix() {
	m := mapify.Map().Set("firstname", "Jane", "partner", struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		32, "Tim",
	})
	m = m.With(mapifier.WithStringPrefix("birth"))("Place", "Map York", "Date", "07/10/1984")

	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"birthDate":"07/10/1984","birthPlace":"Map York","firstname":"Jane","partner":{"age":32,"firstname":"Tim"}}
}

func ExampleArrayIndexes() {
	initialValue := map[string]any{
		"siblings": []struct {
			Age       int32  `structs:"age"`
			Firstname string `structs:"firstname"`
		}{
			{
				Age:       33,
				Firstname: "Tim",
			},
			{
				Age:       31,
				Firstname: "John",
			},
		},
	}
	m := mapify.LoadMap(initialValue)
	m = m.Set("siblings[1].age", 20)
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"siblings":[{"age":33,"firstname":"Tim"},{"age":20,"firstname":"John"}]}
}

func ExampleRootArrayIndexes() {
	m := mapify.Array()
	m = m.Set("[1].age", 20)
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [null,{"age":20}]
}

func ExampleRootArrayOfStructsIndexes() {
	m := mapify.Array()
	m = m.Set("[1]", struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		Age:       33,
		Firstname: "Tim",
	})
	m.Set("[1].age", 23, "[1].lastname", "Doe", "[0].siblings", []string{
		"John", "Jane",
	})
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [{"siblings":["John","Jane"]},{"age":23,"firstname":"Tim","lastname":"Doe"}]
}

func ExampleRootArrayWithSubArrays() {
	var inputValues []any
	inputValues = append(inputValues, struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		Age:       33,
		Firstname: "Tim",
	})
	m := mapify.LoadArray(inputValues)

	m.Set("[0].siblings", []string{
		"John", "Jane",
	})
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [{"age":33,"firstname":"Tim","siblings":["John","Jane"]}]
}

func ExampleRootArrayWithSubArraysAndOverrideTypes() {
	var inputValues []any
	inputValues = append(inputValues, struct {
		Age       int32  `structs:"age"`
		Firstname string `structs:"firstname"`
	}{
		Age:       33,
		Firstname: "Tim",
	})
	m := mapify.LoadArray(inputValues)

	m.Set("[0]", []string{
		"John", "Jane",
	})
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [["John","Jane"]]
}
