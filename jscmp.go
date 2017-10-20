// Package jscmp is a utils to compare two object like js rule
package jscmp

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const (
	zero = json.Number("0")
	one  = json.Number("1")
)

func cmpInt(x, y int64) int {
	if x > y {
		return 1
	} else if x == y {
		return 0
	}
	return -1
}
func cmpFloat(x, y float64) int {
	if x > y {
		return 1
	} else if x == y {
		return 0
	}
	return -1
}

// cmpIntFloat without type convert
// x > y , it will return 1
// x == y, it will return 0
// x < y , if will return -1
func cmpIntFloat(x int64, y float64) int {
	if x == 0 || y == 0 {
		if x == 0 && y == 0 {
			// x == 0, y == 0
			return 0
		}
		if x == 0 {
			if y > 0 {
				return -1
			}
			// x == 0, y > 0
			return 1
		}
		// y == 0
		if x > 0 {
			// x > 0, y == 0
			return 1
		}
		// x < 0, y == 0
		return -1
	} else if x < 0 && y > 0 {
		return -1
	} else if x > 0 && y < 0 {
		return 1
	}
	n := 1
	if x < 0 {
		x = x * -1
		y = math.Abs(y)
		n = -1
	}
	return n * (strings.Compare(fmt.Sprint(x), fmt.Sprint(y)))
}
func isNull(v interface{}) bool {
	return !reflect.ValueOf(v).IsValid()
}

// Equals will return true if left == right like js rule
// but undefined not support in golang
// also not support object
// 0 == null ==> false
// +0 == -0 ==> true
// 1 == 1 ==> true
// 1 == '1' ==> true
// 1 == true ==> true
// 0 == false ==> true
// 0 == undefined ==> false
// null == undefined
func Equals(left, right interface{}) bool {
	if _, ok := left.(json.Number); ok {
		// left is number
		if lefti, ok := parseInt(left); ok {
			// left is int
			if righti, ok := parseInt(right); ok {
				// all int
				return lefti == righti
			}
			if rightf, ok := parseFloat(right); ok {
				// left int, right float
				return cmpIntFloat(lefti, rightf) == 0
			}
			if t, ok := right.(bool); ok {
				if t {
					return lefti == 1
				}
				return lefti == 0
			}
			// eq to null will always return false
			// if right == nil {
			// 	return lefti == 0
			// }
			return false
		}
		if leftf, ok := parseFloat(left); ok {
			if righti, ok := parseInt(right); ok {
				return cmpIntFloat(righti, leftf) == 0
			}
			if rightf, ok := parseFloat(right); ok {
				return leftf == rightf
			}
			if t, ok := right.(bool); ok {
				if t {
					return leftf == 1.0
				}
				return leftf == 0.0
			}
			return false
		}
		// it would not arrive here
		return false
	}
	if _, ok := right.(json.Number); ok {
		return Equals(right, left)
	}
	if t, ok := left.(bool); ok {
		if t {
			return Equals(one, right)
		}
		return Equals(zero, right)
	}
	// left is null and right is null
	if isNull(left) && isNull(right) {
		return true
	}
	// nil map != nil map in js
	// if left ref to the same object with right return true
	if canGetPointer(left) && canGetPointer(right) {
		if reflect.ValueOf(left).Pointer() == reflect.ValueOf(right).Pointer() {
			return true
		}
	}
	if _, ok := parseInt(left); ok {
		return Equals(json.Number(fmt.Sprint(left)), right)
	} else if _, ok := parseFloat(left); ok {
		return Equals(json.Number(fmt.Sprint(left)), right)
	} else if _, ok := parseInt(right); ok {
		return Equals(fmt.Sprint(right), left)
	} else if _, ok := parseFloat(right); ok {
		return Equals(fmt.Sprint(right), left)
	}
	return false
}

// cmp will return true if left > right
// not support object
func cmp(left, right interface{}) bool {
	_, leftNOk := left.(json.Number)
	_, rightNOk := right.(json.Number)
	if leftNOk || rightNOk {
		// at least one of json.Number
		{
			// try treat left as number
			if lefti, ok := parseInt(left); ok {
				// left is int
				if righti, ok := parseInt(right); ok {
					// all int
					return cmpInt(lefti, righti) == 1
				}
				if rightf, ok := parseFloat(right); ok {
					// left int, right float
					return cmpIntFloat(lefti, rightf) == 1
				}
				// right is bool
				if t, ok := right.(bool); ok {
					if t {
						return cmpInt(lefti, 1) == 1
					}
					return cmpInt(lefti, 0) == 1
				}
				// right is null
				if isNull(right) {
					return cmpInt(lefti, 0) == 1
				}
				// left is integer, but right is not number, bool and null
				return false
			}
			// left is not int, try it as float
			if leftf, ok := parseFloat(left); ok {
				// left is float
				if righti, ok := parseInt(right); ok {
					// left float, right int
					return cmpIntFloat(righti, leftf) == -1
				}
				if rightf, ok := parseFloat(right); ok {
					// all float
					return cmpFloat(leftf, rightf) == 1
				}
				// right is bool
				if t, ok := right.(bool); ok {
					if t {
						return cmpFloat(leftf, 1) == 1
					}
					return cmpFloat(leftf, 0) == 1
				}
				// right is null
				if isNull(right) {
					return cmpFloat(leftf, 0) == 1
				}
				return false
			}
		}
		// try left as number failed, and try right as number
		return !cmp(right, left)
	}

	if _, ok := parseInt(left); ok {
		return cmp(json.Number(fmt.Sprint(left)), right)
	} else if _, ok := parseFloat(left); ok {
		return cmp(json.Number(fmt.Sprint(left)), right)
	} else if _, ok := parseInt(right); ok {
		return cmp(left, json.Number(fmt.Sprint(right)))
	} else if _, ok := parseFloat(right); ok {
		return cmp(left, json.Number(fmt.Sprint(right)))
	} else if !reflect.ValueOf(left).IsValid() {
		return cmp(zero, right)
	} else if !reflect.ValueOf(right).IsValid() {
		return cmp(left, zero)
	}
	// left and right both not number
	if t, ok := left.(bool); ok {
		// left is bool
		if t {
			return cmp(one, right)
		}
		return cmp(zero, right)
	}
	if t, ok := right.(bool); ok {
		// right is bool
		if t {
			return cmp(left, one)
		}
		return cmp(left, zero)
	}
	if isNull(left) {
		return cmp(zero, right)
	}
	if isNull(right) {
		return cmp(left, zero)
	}
	return false
}
func parseInt(i interface{}) (int64, bool) {
	s := fmt.Sprint(i)
	if s == "" {
		return 0, true
	}
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return res, true
}
func parseFloat(i interface{}) (float64, bool) {
	res, err := strconv.ParseFloat(fmt.Sprint(i), 64)
	if err != nil {
		return 0, false
	}
	return res, true
}

// GT is >
func GT(x, y interface{}) bool {
	return cmp(x, y)
}

// GTE is >=
func GTE(x, y interface{}) bool {
	if cmp(x, y) {
		return true
	}
	if !reflect.ValueOf(x).IsValid() {
		return Equals(zero, y)
	} else if !reflect.ValueOf(y).IsValid() {
		return Equals(x, zero)
	}
	return Equals(x, y)
}

// LT is <
func LT(x, y interface{}) bool {
	if cmp(x, y) {
		return false
	}
	return !Equals(x, y)
}

// LTE is <=
func LTE(x, y interface{}) bool {
	return !cmp(x, y)
}

// Not is !
func Not(x interface{}) bool {
	if isNull(x) || x == 0 || x == false || x == "" {
		return true
	}
	if n, ok := parseFloat(x); ok {
		if n == 0 {
			return true
		}
	}
	return false
}

func checkComparable(i interface{}) bool {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Map, reflect.Array, reflect.Func, reflect.Chan, reflect.Struct, reflect.Slice, reflect.UnsafePointer, reflect.Interface:
		return false
	}
	return true
}

func canGetPointer(i interface{}) bool {
	if reflect.ValueOf(i).IsValid() {
		return false
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

// StrictEquals like js ===
func StrictEquals(x, y interface{}) bool {
	if isNull(x) && !isNull(y) {
		return false
	} else if !isNull(x) && isNull(y) {
		return false
	}
	if canGetPointer(x) && canGetPointer(y) {
		if reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer() {
			return true
		}
	}

	if !checkComparable(x) || !checkComparable(y) {
		return false
	}

	if x == y {
		return true
	}
	var (
		xi, yi *int64
		xf, yf *float64
	)
	switch reflect.TypeOf(x).Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// TODO(cloverstd): support uint64
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		xi = new(int64)
		*xi = reflect.ValueOf(x).Int()
	case reflect.Float32, reflect.Float64:
		xf = new(float64)
		*xf = reflect.ValueOf(x).Float()
	default:
		return false
	}
	switch reflect.TypeOf(y).Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// TODO(cloverstd): support uint64
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		yi = new(int64)
		*yi = reflect.ValueOf(y).Int()
	case reflect.Float32, reflect.Float64:
		yf = new(float64)
		*yf = reflect.ValueOf(y).Float()
	default:
		return false
	}

	if xi != nil && yi != nil {
		return cmpInt(*xi, *yi) == 0
	} else if xf != nil && yf != nil {
		return cmpFloat(*xf, *yf) == 0
	} else if xi != nil && yf != nil {
		return cmpIntFloat(*xi, *yf) == 0
	} else if xf != nil && yi != nil {
		return cmpIntFloat(*yi, *xf) == 0
	}
	return false
}
