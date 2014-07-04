package devatlas

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
)

const API_VERSION = "6"

// OpenFile opens a data file
func OpenFile(name string) (*Atlas, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return Open(file)
}

// Open opens a reader
func Open(r io.Reader) (*Atlas, error) {
	var atlas Atlas
	err := json.NewDecoder(r).Decode(&atlas)
	if err != nil {
		return nil, err
	}

	exprs, _ := atlas.Expressions[API_VERSION]
	atlas.regexes = make([]*regexp.Regexp, len(exprs))
	for i, expr := range exprs {
		atlas.regexes[i] = regexp.MustCompile(expr)
	}

	return &atlas, nil
}

// Find looks up the user agent string
func (db *Atlas) Find(ua string) map[string]interface{} {
	acc := make(map[int]int)
	db.Tree.traverse(ua, db.regexes, acc)

	attrs := make(map[string]interface{}, len(acc))
	for pi, vi := range acc {
		prop := db.Properties[pi]
		attrs[prop.Name] = prop.Convert(db.Values[vi])
	}
	return attrs
}

func (n *Node) traverse(ua string, regexes []*regexp.Regexp, acc map[int]int) {
	for pi, vi := range n.Data {
		acc[pi] = vi
	}
	if len(n.Children) < 1 {
		return
	}

	for _, pos := range n.Regexes {
		ua = regexes[pos].ReplaceAllString(ua, "")
	}

	for i := 1; i <= len(ua); i++ {
		if child, ok := n.Children[ua[:i]]; ok {
			child.traverse(ua[i:], regexes, acc)
		}
	}
}
