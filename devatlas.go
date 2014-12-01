package devatlas

import (
	"encoding/json"
	"io"
	"os"
)

// API version
const API_VERSION = "6"

// Property kinds
const (
	KIND_BOOL uint8 = 1 << iota
	KIND_INT
	KIND_STRING
)

type DB struct{ rawData }

// OpenFile opens a data file
func OpenFile(name string) (*DB, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return Open(file)
}

// Open opens a reader
func Open(r io.Reader) (*DB, error) {
	var raw rawData
	err := json.NewDecoder(r).Decode(&raw)
	if err != nil {
		return nil, err
	}
	return &DB{raw}, nil
}

// Find looks up the user agent string
func (db *DB) Find(ua string) map[string]interface{} {
	acc := make(indexMap, 100)
	db.Tree.traverse(ua, db.Regexp, acc)

	overrides := db.UAR.update(ua, acc)
	attrs := make(map[string]interface{}, len(acc))
	for pi, vi := range acc {
		prop := db.Properties[pi]
		if val, ok := overrides[pi]; ok {
			attrs[prop.Name] = prop.Convert(val)
		} else {
			attrs[prop.Name] = prop.Convert(db.Values[vi])
		}
	}
	return attrs
}
