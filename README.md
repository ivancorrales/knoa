[![GitHub Release](https://img.shields.io/github/v/release/ivancorrales/knoa)](https://github.com/ivancorrales/knoa/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/ivancorrales/knoa.svg)](https://pkg.go.dev/github.com/ivancorrales/knoa)
[![go.mod](https://img.shields.io/github/go-mod/go-version/ivancorrales/knoa)](go.mod)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://img.shields.io/github/license/ivancorrales/knoa)
[![Build Status](https://img.shields.io/github/actions/workflow/status/ivancorrales/knoa/build.yml?branch=main)](https://github.com/ivancorrales/knoa/actions?query=workflow%3ABuild+branch%3Amain)
[![CodeQL](https://github.com/ivancorrales/knoa/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/ivancorrales/knoa/actions/workflows/codeql.yml)

# Knoa

The `swiss knife` to deal with the hassle of `unstructured data`.

## History and project status

This module is already `ready-for-production`.

## Knoa  Highlights

* **Easy integration**: It's straightforward to be integrated with your current developments. 

## Installation

Use go get to retrieve the library to add it to your GOPATH workspace, or project's Go module dependencies.

```bash
go get -u github.com/ivancorrales/knoa
```

To update the library use go get -u to retrieve the latest version of it.

```bash
go get -u github.com/ivancorrales/knoa
```

You could specify a concrete version of this module as It's shown on the below. Replace x.y.z by the desired version.

```bash
module github.com/<org>/<repository>
require ( 
  github.com/ivancorrales/knoa vX.Y.Z
)
```

## Getting started

### Pre-requisites

* Go 1.19+

### Examples

Check the folder `examples`

```go
package main

import (
	"fmt"
	"strings"

	"github.com/ivancorrales/knoa"
	"github.com/ivancorrales/knoa/mapifier"
)

// Basic showcase
func ExampleFromScratch() {
	out := knoa.Map().Set("firstname", "Jane").JSON()
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
	knoa.Map(mapifier.WithStrictMode(true)).Set("firstname.$", "Jane").JSON()
	// Output:
	// invalid Path  'firstname.$'. Path doesn't match defined format
}

// Ignore those attributes that don't match the provided format
func ExampleFromScratchWithAttributeNameFormat() {
	out := knoa.Map(mapifier.WithAttributeNameFormat("person-(.*)")).Set("person-firstname", "Jane", "lastname", "Doe").JSON()
	fmt.Println(out)
	// Output:
	// {"person-firstname":"Jane"}
}

// Set arrays attributes
func ExampleSetArrayChildren() {
	out := knoa.Map().Set("firstname", "Jane", "siblings", []string{"Tim", "Janet"}).JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","siblings":["Tim","Janet"]}
}

// Set arrays attributes
func ExampleSetArrayChildrenV2() {
	out := knoa.Map().Set("firstname", "Jane", "languages.native", []string{"English", "Irish"}, "languages.learning", []string{"Italian"}).JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Jane","languages":{"learning":["Italian"],"native":["English","Irish"]}}
}

// Set several times
func ExampleMultipleSets() {
	m := knoa.Map().Set("firstname", "Jane")
	m = m.Set("lastname", "Doe")
	m = m.Set("firstname", "Tim")
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"firstname":"Tim","lastname":"Doe"}
}

// Set complex structures
func ExampleSetComplexStructures() {
	out := knoa.Map().Set("firstname", "Jane", "partner", struct {
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
	m := knoa.Map().Set("firstname", "Jane", "partner", struct {
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
	m := knoa.Map().Set("firstname", "Jane", "partner", struct {
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
	m := knoa.Map().Set("firstname", "Jane", "partner", struct {
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
	m := knoa.LoadMap(initialValue)
	m = m.Set("siblings[1].age", 20)
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// {"siblings":[{"age":33,"firstname":"Tim"},{"age":20,"firstname":"John"}]}
}

func ExampleRootArrayIndexes() {
	m := knoa.Array()
	m = m.Set("[1].age", 20)
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [null,{"age":20}]
}

func ExampleRootArrayOfStructsIndexes() {
	m := knoa.Array()
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
	m := knoa.LoadArray(inputValues)

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
	m := knoa.LoadArray(inputValues)

	m.Set("[0]", []string{
		"John", "Jane",
	})
	out := m.JSON()
	fmt.Println(out)
	// Output:
	// [["John","Jane"]]
}
```

### Contributing

See the [contributing](https://github.com/ivancorrales/knoa/blob/main/CONTRIBUTING.md) documentation.


