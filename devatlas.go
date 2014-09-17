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

	if err = atlas.parse(); err != nil {
		return nil, err
	}
	return &atlas, nil
}

// Find looks up the user agent string
func (db *Atlas) Find(ua string) map[string]interface{} {
	acc := make(Indices)
	db.Tree.traverse(ua, db.regexes, acc)

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

func (n *Node) traverse(ua string, regexes []*regexp.Regexp, acc Indices) {
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

func (u *UAR) update(ua string, acc Indices) map[int]interface{} {
	// Check if any of the pre-resolved properties require skipping UAR traversal
	for _, pi := range u.Skip {
		if acc[pi] > 0 {
			return nil
		}
	}

	res := make(map[int]interface{})
	for _, rg := range u.RuleGroups {
		for _, rule := range rg.Rules(ua, u.Regexes, acc) {
			if rule.Vi != nil {
				acc[rule.Pi] = *rule.Vi
			} else if m := u.Regexes[rule.Ri].FindStringSubmatch(ua); rule.M < len(m) {
				res[rule.Pi] = m[rule.M]
			}
		}
	}
	return res
}
