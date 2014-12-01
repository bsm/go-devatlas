package devatlas

import "regexp"

// Tree nodes
type treeNode struct {
	Children    map[string]treeNode `json:"c,omitempty"`
	Data        indexMap            `json:"d,omitempty"`
	Corrections []int               `json:"r,omitempty"`

	// tlns []int // index of child token lengths
}

// index child token lengths
// func (n *treeNode) index() {
// 	n.tlns = make([]int, 0, len(n.Children))
// 	iset := make(map[int]bool, len(n.Children))
// 	for str, _ := range n.Children {
// 		iset[len(str)] = true
// 	}
// 	for num, _ := range iset {
// 		n.tlns = append(n.tlns, num)
// 	}
// }

func (n *treeNode) traverse(ua string, regexes []*regexp.Regexp, acc indexMap) {
	// Merge node attributes into acc
	for pid, vid := range n.Data {
		acc[pid] = vid
	}

	// Done, if there are no children
	if len(n.Children) < 1 {
		return
	}

	// Apply corrections
	for _, pos := range n.Corrections {
		ua = regexes[pos].ReplaceAllString(ua, "")
	}

	// Find next child
	for i := 1; i <= len(ua); i++ {
		if child, ok := n.Children[ua[:i]]; ok {
			child.traverse(ua[i:], regexes, acc)
			return
		}
	}
}
