package jscmp

import "encoding/json"

// Number represents a number could be compared
type Number interface {
	Int64() (int64, error)
	Float64() (float64, error)
	String() string
}

var _ Number = (*json.Number)(nil)
