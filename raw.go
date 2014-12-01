package devatlas

import (
	"encoding/json"
	"strconv"

	"github.com/martinolsen/go-pcre/regexp"
)

// Raw data format, as stored in the data file
type rawData struct {
	Meta       Meta         `json:"$"`
	Tree       treeNode     `json:"t"`
	Values     mixedSlice   `json:"v"`
	Properties []Property   `json:"p"`
	Regexp     regexpSliceV `json:"r"`
	UAR        UAR          `json:"uar"`
}

// Meta data, simple values
type Meta struct {
	Ver, Rev  string
	Trademark string
	Gen       string
	Utc       int64
}

// Index maps map propertyIDs to valueIDs
type indexMap map[int]int

// UnmarshalJSON custom implementation
//   {"7":12314,"8":63662}
func (i *indexMap) UnmarshalJSON(b []byte) error {
	var kv map[string]int
	if err := json.Unmarshal(b, &kv); err != nil {
		return err
	}

	res := make(indexMap, len(kv))
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

// List of versioned regular expressions
//   {
//     "5":["A","B"],
//     "6":["C","D"]
//   }
type regexpSliceV []*regexp.Regexp

func (s *regexpSliceV) UnmarshalJSON(b []byte) error {
	var raw map[string][]string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	plain, _ := raw[API_VERSION]
	slice := make(regexpSliceV, len(plain))
	for i, expr := range plain {
		x, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		slice[i] = x
	}

	*s = slice
	return nil
}

// List of regular expressions with defaults and version overrides
//   {
//     "D":["A","B"],
//     "6":{"1":"C"}
//   }
type regexpSliceO []*regexp.Regexp

func (s *regexpSliceO) UnmarshalJSON(b []byte) error {
	var raw struct {
		Defaults  []string          `json:"d"`
		Overrides map[string]string `json:"6"` // set to API_VERSION
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	slice := make(regexpSliceO, 0, len(raw.Defaults))
	for _, expr := range raw.Defaults {
		x, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		slice = append(slice, x)
	}

	for s, expr := range raw.Overrides {
		x, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		slice[i] = x
	}

	*s = slice
	return nil
}

// mixedSlice are an slice of strings and uint32s
//    ["A", 1, "B", 2]
type mixedSlice []interface{}

// UnmarshalJSON custom implementation
func (v *mixedSlice) UnmarshalJSON(b []byte) error {
	var values []interface{}
	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}

	for i, item := range values {
		switch t := item.(type) {
		case int:
			values[i] = int(t)
		case int8:
			values[i] = int(t)
		case int16:
			values[i] = int(t)
		case int32:
			values[i] = int(t)
		case int64:
			values[i] = int(t)
		case uint:
			values[i] = int(t)
		case uint8:
			values[i] = int(t)
		case uint16:
			values[i] = int(t)
		case uint32:
			values[i] = int(t)
		case uint64:
			values[i] = int(t)
		case float32:
			values[i] = int(t)
		case float64:
			values[i] = int(t)
		}
	}

	*v = mixedSlice(values)
	return nil
}
