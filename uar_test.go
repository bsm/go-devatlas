package devatlas

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAR", func() {

	It("should parse", func() {
		Expect(testDB.UAR).NotTo(BeNil())
		Expect(testDB.UAR.Skip).NotTo(BeEmpty())
		Expect(testDB.UAR.Regexp).NotTo(BeEmpty())
		Expect(testDB.UAR.RuleGroups).NotTo(BeEmpty())

		group := testDB.UAR.RuleGroups[0]
		Expect(group.Sets).NotTo(BeEmpty())
		Expect(group.Reqs).NotTo(BeEmpty())
	})

})

var _ = Describe("ruleGroup", func() {
	r1 := rule{Pi: 123, Ri: 0, M: 2}
	r2 := rule{Pi: 162, Vi: iptr(1080), Ri: 1, M: 1}
	testCases := []struct {
		desc string
		rg   ruleGroup
		acc  indexMap
		exp  []rule
	}{
		{`no requirements - no rules`,
			ruleGroup{Reqs: indexMap{}, Sets: []ruleSet{ruleSet{Rules: []rule{r1}}}},
			indexMap{1: 10}, nil},
		{`requirements not met - no rules`,
			ruleGroup{Reqs: indexMap{1: 11}, Sets: []ruleSet{ruleSet{Rules: []rule{r1}}}},
			indexMap{1: 10}, nil},
		{`not all requirements met - no rules`,
			ruleGroup{Reqs: indexMap{1: 10, 2: 20}, Sets: []ruleSet{ruleSet{Rules: []rule{r1}}}},
			indexMap{1: 10}, nil},
		{`only one rule set - use the rules`,
			ruleGroup{Reqs: indexMap{1: 10}, Sets: []ruleSet{ruleSet{Rules: []rule{r1}}}},
			indexMap{1: 10}, []rule{r1}},
		{`no rule set matches - no rules`,
			ruleGroup{Reqs: indexMap{1: 10}, Sets: []ruleSet{ruleSet{Rules: []rule{r1}}, ruleSet{Rules: []rule{r2}}}},
			indexMap{1: 10}, []rule{}},
		{`a rule set matches - use the rules`,
			ruleGroup{Reqs: indexMap{1: 10}, Sets: []ruleSet{ruleSet{Mi: iptr(3), Rules: []rule{r1}}, ruleSet{Mi: iptr(2), Rules: []rule{r2}}}},
			indexMap{1: 10}, []rule{r2}},
		{`multiple rule set match - merge the rules`,
			ruleGroup{Reqs: indexMap{1: 10}, Sets: []ruleSet{ruleSet{Mi: iptr(2), Rules: []rule{r1}}, ruleSet{Mi: iptr(2), Rules: []rule{r2}}}},
			indexMap{1: 10}, []rule{r1, r2}},
	}

	It("should extract rules", func() {
		for _, tc := range testCases {
			rules := tc.rg.matchRules(`Mozilla/5.0 (iPod touch; CPU iPhone OS 7_0_4 like Mac OS X)`, testDB.UAR.Regexp, tc.acc)
			if tc.exp == nil {
				Expect(rules).To(BeNil(), tc.desc)
			} else {
				Expect(rules).To(Equal(tc.exp), tc.desc)
			}
		}
	})

})

func iptr(i int) *int { return &i }
