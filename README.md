[![Build Status](https://travis-ci.org/cloverstd/jscmp.svg?branch=master)](https://travis-ci.org/cloverstd/jscmp) [![Go Report Card](https://goreportcard.com/badge/github.com/cloverstd/jscmp)](https://goreportcard.com/report/github.com/cloverstd/jscmp) [![codecov](https://codecov.io/gh/cloverstd/jscmp/branch/master/graph/badge.svg)](https://codecov.io/gh/cloverstd/jscmp)
## Introduction

The jscmp package can compare two objects like JavaScript compare rules.

## Installation and usage

```bash
go get github.com/cloverstd/jscmp
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/cloverstd/jscmp"
)

func Example() {
	var x, y interface{}
	json.Unmarshal([]byte("true"), &x)
	json.Unmarshal([]byte("1"), &y)

	// true == 1
	fmt.Println(jscmp.Equals(x, y))
	// true === 1
	fmt.Println(jscmp.StrictEquals(x, y))

	json.Unmarshal([]byte("10"), &x)
	json.Unmarshal([]byte(`"1"`), &y)
	// 10 > '1'
	fmt.Println(jscmp.GTE(x, y))
	// 10 > true
	fmt.Println(jscmp.GTE(10, true))

	json.Unmarshal([]byte("-10"), &x)
	json.Unmarshal([]byte(`"-1.00"`), &y)
	// -10 > "-1.00"
	fmt.Println(jscmp.GT(x, y))
	// -10 < "-1.00"
	fmt.Println(jscmp.LT(x, y))
}
```

This example will generate the following output:
```bash
true
false
true
true
false
true
```

## NOT SUPPORT

* `[]string{}` is the same object in Golang, if you use `jscmp.Equals([]string{}, []string{})`, it will return true. Is js `{} == {}` will return false.
* object cmp not support now, cmp two objects with jscmp, if the objects both reference to the same object, you will get true when cmp them.
