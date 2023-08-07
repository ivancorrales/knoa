[![GitHub Release](https://img.shields.io/github/v/release/ivancorrales/knoa)](https://github.com/ivancorrales/knoa/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/ivancorrales/knoa.svg)](https://pkg.go.dev/github.com/ivancorrales/knoa)
[![go.mod](https://img.shields.io/github/go-mod/go-version/ivancorrales/knoa)](go.mod)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://img.shields.io/github/license/ivancorrales/knoa)
[![Build Status](https://img.shields.io/github/actions/workflow/status/ivancorrales/knoa/build.yml?branch=main)](https://github.com/ivancorrales/knoa/actions?query=workflow%3ABuild+branch%3Amain)
[![CodeQL](https://github.com/ivancorrales/knoa/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/ivancorrales/knoa/actions/workflows/codeql.yml)

# Knoa

The `swiss knife` to deal with the hassle of `unstructured data`.

## Getting started

### Installation

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

### Pre-requisites

* Go 1.19+

### Examples

Check the following examples


**Create and modify an array dynamically**

```go
k := knoa.Array().Set("[1]", []string{"red", "blue", "green"}, "[2].firstname", "John")
fmt.Println(k.String())
// [null,["red","blue","green"],{"firstname":"John"}]

k.Set("[0]", struct {
    FullName  string `structs:"fullName"`
    RoleLevel int    `structs:"roleLevel"`
}{
    "Senior Developer",
    3,
})
fmt.Println(k.String())
// [{"fullName":"Senior Developer","roleLevel":3},["red","blue","green"],{"firstname":"John"}]

k.Set("[0].enabled", false, "[2].firstname", "Jane")
fmt.Println(k.String())
// [{"enabled":false,"fullName":"Senior Developer","roleLevel":3},["red","blue","green"],{"firstname":"Jane"}]
```

**Create and modify a map dynamically**

```go
k := knoa.Map().Set("firstname", "John", "age", 20)
fmt.Println(k.String())
// {"age":20,"firstname":"John"}

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
// {"age":20,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":40,"firstname":"Bob"}]}

k.Set("age", 23, "siblings[1].age", 39)
fmt.Println(k.String())
// {"age":23,"firstname":"John","siblings":[{"age":29,"firstname":"Tim"},{"age":39,"firstname":"Bob"}]}
```



**Decode an unstructured data into a struct and an array of structs**
```go

type Person struct {
    Firstname string `struct:"firstname"`
    Age       int    `struct:"age"`
}



k:= knoa.Load[[]any]([]any{
    Person{
        Firstname: "Jane",
        Age:       20,
	},
})
k.Set("[1]", Person{
        Firstname: "Bob",
        Age:       23,
    },"[2].firstname", "John")

var output []Person
k.To(&output)
```

**Decode an unstructured data into a struct**
```go

type Person struct {
    Firstname string   `structs:"firstname"`
    Age       int      `structs:"age"`
    Siblings  []Person `structs:"siblings,omitempty"`
}

k := knoa.Map().Set("firstname", "John", "age", 20)
	
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

k.Set("age", 23, "siblings[1].age", 39)

var person Person
k.To(&person)
```


Additionally, we encourege to have a look at folder `examples` to get a better understanding on how `knoa` works.

### Contributing

See the [contributing](https://github.com/ivancorrales/knoa/blob/main/CONTRIBUTING.md) documentation.


