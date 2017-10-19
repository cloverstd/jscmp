package jscmp_test

import (
	"encoding/json"
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
	// Output:
	// true
	// false
	// true
	// true
	// false
	// true
}
