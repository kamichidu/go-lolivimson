package lolivimson

import (
	"fmt"
	ref "reflect"
	"sort"
	"strings"
)

const (
	INT32_MAX = 2147483647
	INT32_MIN = -2147483648
)

var SortDictionaryKey bool = false

type Encoder struct {
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (self *Encoder) Marshal(v interface{}) ([]byte, error) {
	return self.encodeValue(ref.ValueOf(v))
}

func (self *Encoder) encodeValue(v ref.Value) ([]byte, error) {
	t := v.Type()
	switch t.Kind() {
	case ref.Ptr:
		return self.encodeValue(v.Elem())
	case ref.String:
		return self.encodeString(v)
	case ref.Bool:
		return self.encodeBool(v)
	case ref.Int, ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Uint, ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64:
		return self.encodeNumber(v)
	case ref.Float32, ref.Float64:
		return self.encodeFloat(v)
	case ref.Array, ref.Slice:
		return self.encodeList(v)
	case ref.Map:
		return self.encodeDictionary(v)
	default:
		return nil, fmt.Errorf("Unsupported type: %s", t.String())
	}
}

func (self *Encoder) encodeBool(v ref.Value) ([]byte, error) {
	val := v.Bool()
	if val {
		return []byte("1"), nil
	} else {
		return []byte("0"), nil
	}
}

func (self *Encoder) encodeString(v ref.Value) ([]byte, error) {
	val := v.String()
	return []byte("'" + strings.Replace(val, "'", "''", -1) + "'"), nil
}

func (self *Encoder) encodeNumber(v ref.Value) ([]byte, error) {
	val := v.Int()
	// http://golang.org/ref/spec#Numeric_types
	if val > INT32_MAX {
		return nil, fmt.Errorf("Overflow number value range: %d", val)
	} else if val < INT32_MIN {
		return nil, fmt.Errorf("Underflow number value range: %d", val)
	}

	return []byte(fmt.Sprintf("%d", val)), nil
}

func (self *Encoder) encodeFloat(v ref.Value) ([]byte, error) {
	val := v.Float()

	return []byte(fmt.Sprintf("%.9f", val)), nil
}

func (self *Encoder) encodeList(v ref.Value) ([]byte, error) {
	elements := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		elm, err := self.encodeValue(v.Index(i))
		if err != nil {
			return nil, err
		}
		elements[i] = string(elm)
	}
	return []byte("[" + strings.Join(elements, ",") + "]"), nil
}

type SortableValues []ref.Value

func (self SortableValues) Len() int {
	return len(self)
}
func (self SortableValues) Swap(i int, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self SortableValues) Less(i int, j int) bool {
	return self[i].String() < self[j].String()
}

func (self *Encoder) encodeDictionary(v ref.Value) ([]byte, error) {
	items := make([]string, v.Len())
	keys := v.MapKeys()

	if SortDictionaryKey {
		sort.Sort(SortableValues(keys))
	}

	for i := 0; i < v.Len(); i++ {
		key, err := self.encodeValue(keys[i])
		if err != nil {
			return nil, err
		}
		value, err := self.encodeValue(v.MapIndex(keys[i]))
		if err != nil {
			return nil, err
		}
		items[i] = strings.Join([]string{string(key), string(value)}, ":")
	}
	return []byte("{" + strings.Join(items, ",") + "}"), nil
}
