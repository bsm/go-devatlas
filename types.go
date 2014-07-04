package devatlas

import (
	"encoding/json"
	"regexp"
	"strconv"
)

const (
	KIND_BOOL uint8 = 1 << iota
	KIND_INT
	KIND_STRING
)

// Main Atlas DB
type Atlas struct {
	Meta        *Meta               `json:"$"`
	Tree        *Node               `json:"t"`
	Values      Values              `json:"v"`
	Properties  []Property          `json:"p"`
	Expressions map[string][]string `json:"r"`

	regexes []*regexp.Regexp
}

// Meta data
type Meta struct {
	Ver, Rev  string
	Trademark string
	Gen       string
	Utc       int64
}

// Tree nodes
type Node struct {
	Children map[string]Node `json:"c,omitempty"`
	Data     Indices         `json:"d,omitempty"`
	Regexes  []int           `json:"r,omitempty"`
}

type Values []interface{}

// UnmarshalJSON custom implementation
func (v *Values) UnmarshalJSON(b []byte) error {
	var is []interface{}
	if err := json.Unmarshal(b, &is); err != nil {
		return err
	}

	for i, item := range is {
		switch t := item.(type) {
		case int:
			is[i] = int(t)
		case int8:
			is[i] = int(t)
		case int16:
			is[i] = int(t)
		case int32:
			is[i] = int(t)
		case int64:
			is[i] = int(t)
		case uint:
			is[i] = int(t)
		case uint8:
			is[i] = int(t)
		case uint16:
			is[i] = int(t)
		case uint32:
			is[i] = int(t)
		case uint64:
			is[i] = int(t)
		case float32:
			is[i] = int(t)
		case float64:
			is[i] = int(t)
		}
	}

	*v = Values(is)
	return nil
}

type Property struct {
	Name string
	Kind uint8
}

// Convert converts a value to the `kind`interface
func (p *Property) Convert(v interface{}) interface{} {
	switch p.Kind {
	case KIND_BOOL:
		return v == 1
	}
	return v
}

// UnmarshalJSON custom implementation
func (p *Property) UnmarshalJSON(b []byte) error {
	var id string
	if err := json.Unmarshal(b, &id); err != nil {
		return err
	}

	kind := KIND_STRING
	switch id[0] {
	case 'b':
		kind = KIND_BOOL
	case 'i':
		kind = KIND_INT
	case 's':
	default:
		panic("unknown kind " + id)
	}
	*p = Property{Name: id[1:], Kind: kind}
	return nil
}

// Index map
type Indices map[int]int

// UnmarshalJSON custom implementation
func (i *Indices) UnmarshalJSON(b []byte) error {
	var kv map[string]int
	if err := json.Unmarshal(b, &kv); err != nil {
		return err
	}

	res := make(Indices)
	for str, val := range kv {
		num, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		res[num] = val
	}
	*i = res
	return nil
}
