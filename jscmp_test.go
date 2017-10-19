package jscmp_test

import (
	"encoding/json"
	"testing"

	. "github.com/cloverstd/jscmp"
)

func TestEquals(t *testing.T) {
	if !Equals(+0, -0) {
		t.Error("test +0 == -0 failed")
	}
	if !Equals(1, 1) {
		t.Error("test 1 == 1 failed")
	}

	if !Equals(1, "1") {
		t.Error("test 1 == '1' failed")
	}

	if !Equals(1, true) {
		t.Error("test 1 == true failed")
	}

	if !Equals(0, false) {
		t.Error("test 0 == false failed")
	}

	if !Equals("1", true) {
		t.Error("test '1' == true failed")
	}

	if !Equals("0", false) {
		t.Error("test '0' == false failed")
	}

	if !Equals("", false) {
		t.Error("test '' == false failed")
	}

	if Equals(nil, 0) {
		t.Error("test null != 0 failed")
	}

	if Equals(nil, "0") {
		t.Error("test null != '0' failed")
	}

	if Equals(nil, "") {
		t.Error("test null != '' failed")
	}

	if Equals(map[int]int{}, nil) {
		t.Error("test {} != null failed")
	}

	if Equals(map[int]int{}, map[int]int{}) {
		t.Error("test {} != {} failed")
	}

	m1 := map[int]int{1: 1}
	m2 := m1
	if !Equals(m1, m2) {
		t.Error("test ref {} == {} failed")
	}

	if !Equals(1.0, "1.00") {
		t.Error("test 1.0 == '1.00' failed")
	}

	if !Equals(1.0, 1) {
		t.Error("test 1.0 == 1 failed")
	}
}

func TestGT(t *testing.T) {
	if !GT(2, true) {
		t.Error("test 2 > true failed")
	}

	if !GT(1, false) {
		t.Error("test 1 > false failed")
	}

	if !GT(10, 0) {
		t.Error("test 10 > 0 failed")
	}

	if !GT(10, "-10") {
		t.Error("test 10 > '-10' failed")
	}

	if !GT(-1, -10) {
		t.Error("test -1 > -10 failed")
	}

	if !GT(10, -10) {
		t.Error("test 10 > -10 failed")
	}
}

func TestGTE(t *testing.T) {
	if !GTE(1, true) {
		t.Error("test 1 >= true failed")
	}

	if !GTE(2, true) {
		t.Error("test 2 >= true failed")
	}

	if GTE(2, "true") {
		t.Error("test !(2 >= 'true') failed")
	}

	if !GTE(0, nil) {
		t.Error("test 0 >= nil failed")
	}

	if !GT(10, nil) {
		t.Error("test 10 >= nil failed")
	}

	if !GTE(0, "") {
		t.Error("test 0 >= '' failed")
	}

	if !GTE(10, "") {
		t.Error("test 10 >= '' failed")
	}

	if GTE(10, map[int]int{}) {
		t.Error("test !(10 >= {}) failed")
	}
}

func TestStrictEquals(t *testing.T) {
	if !StrictEquals(1, 1) {
		t.Error("test 1 === 1 failed")
	}

	if StrictEquals(1, "1") {
		t.Error("test 1 === '1' failed")
	}

	if StrictEquals(true, 1) {
		t.Error("test true === 1 failed")
	}

	if StrictEquals(false, 0) {
		t.Error("test false === 0 failed")
	}

	if !StrictEquals(1, 1.0) {
		t.Error("test 1 === 1.0 failed")
	}

	if !StrictEquals(-1.0, -1) {
		t.Error("test -1.0 === -1 failed")
	}

	if !StrictEquals(0, -0) {
		t.Error("test 0 === -0 failed")
	}
	var (
		x, y interface{}
	)
	json.Unmarshal([]byte(`{}`), &x)
	json.Unmarshal([]byte(`{}`), &y)
	if StrictEquals(x, y) {
		t.Error("test !([] === []) failed")
	}
}