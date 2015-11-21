package test

import (
	v "github.com/kamichidu/go-lolivimson"
	"reflect"
	"testing"
)

func init() {
	v.SortDictionaryKey = true
}

func ok(t *testing.T, val interface{}, expect string) {
	enc := v.NewEncoder()
	ret, err := enc.Marshal(val)
	if err != nil {
		t.Errorf("Passed %v then got %v", val, err)
	} else if !reflect.DeepEqual(ret, []byte(expect)) {
		t.Errorf("Expected %s, got %s", expect, ret)
	}
}

func TestEncoder_Marshal_bool(t *testing.T) {
	ok(t, true, "1")
	ok(t, false, "0")

	var val bool
	val = true
	ok(t, &val, "1")

	val = false
	ok(t, &val, "0")
}

func TestEncoder_Marshal_string(t *testing.T) {
	ok(t, "hello", "'hello'")
	ok(t, "はろー'''てすと'''", "'はろー''''''てすと'''''''")
}

func TestEncoder_Marshal_number(t *testing.T) {
	ok(t, 0, "0")
	ok(t, 1, "1")
	ok(t, -1, "-1")
}

func TestEncoder_Marshal_float(t *testing.T) {
	ok(t, 0.0, "0.000000000")
	ok(t, 1e3, "1000.000000000")
	ok(t, -1e3, "-1000.000000000")
	ok(t, 3.14, "3.140000000")
}

func TestEncoder_Marshal_list(t *testing.T) {
	ok(t, []string{}, "[]")
	ok(t, []string{"a", "b", "c"}, "['a','b','c']")
}

func TestEncoder_Marshal_dictionary(t *testing.T) {
	var m map[string]int

	m = make(map[string]int)
	ok(t, m, "{}")

	m = make(map[string]int)
	m["one"] = 1
	m["two"] = 2
	m["three"] = 3
	ok(t, m, "{'one':1,'three':3,'two':2}")
}
