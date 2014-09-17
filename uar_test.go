package devatlas

import (
	"encoding/json"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAR", func() {

	It("should parse", func() {
		Expect(testDB.UAR).NotTo(BeNil())
		Expect(testDB.UAR.Skip).NotTo(BeEmpty())
		Expect(testDB.UAR.Regexes).NotTo(BeEmpty())
		Expect(testDB.UAR.RuleGroups).NotTo(BeEmpty())

		group := testDB.UAR.RuleGroups[0]
		Expect(group.Sets).NotTo(BeEmpty())
		Expect(group.Reqs).NotTo(BeEmpty())
	})

	Describe("RuleGroup", func() {
		r1 := Rule{Pi: 123, Ri: 0, M: 2}
		r2 := Rule{Pi: 162, Vi: iptr(1080), Ri: 1, M: 1}
		testCases := []struct {
			desc string
			rg   RuleGroup
			acc  Indices
			exp  []Rule
		}{
			{`no requirements - no rules`,
				RuleGroup{Reqs: Indices{}, Sets: []RuleSet{RuleSet{Rules: []Rule{r1}}}},
				Indices{1: 10}, nil},
			{`requirements not met - no rules`,
				RuleGroup{Reqs: Indices{1: 11}, Sets: []RuleSet{RuleSet{Rules: []Rule{r1}}}},
				Indices{1: 10}, nil},
			{`not all requirements met - no rules`,
				RuleGroup{Reqs: Indices{1: 10, 2: 20}, Sets: []RuleSet{RuleSet{Rules: []Rule{r1}}}},
				Indices{1: 10}, nil},
			{`only one rule set - use the rules`,
				RuleGroup{Reqs: Indices{1: 10}, Sets: []RuleSet{RuleSet{Rules: []Rule{r1}}}},
				Indices{1: 10}, []Rule{r1}},
			{`no rule set matches - no rules`,
				RuleGroup{Reqs: Indices{1: 10}, Sets: []RuleSet{RuleSet{Rules: []Rule{r1}}, RuleSet{Rules: []Rule{r2}}}},
				Indices{1: 10}, []Rule{}},
			{`a rule set matches - use the rules`,
				RuleGroup{Reqs: Indices{1: 10}, Sets: []RuleSet{RuleSet{Mi: iptr(3), Rules: []Rule{r1}}, RuleSet{Mi: iptr(2), Rules: []Rule{r2}}}},
				Indices{1: 10}, []Rule{r2}},
			{`multiple rule set match - merge the rules`,
				RuleGroup{Reqs: Indices{1: 10}, Sets: []RuleSet{RuleSet{Mi: iptr(2), Rules: []Rule{r1}}, RuleSet{Mi: iptr(2), Rules: []Rule{r2}}}},
				Indices{1: 10}, []Rule{r1, r2}},
		}

		It("should extract rules", func() {
			for _, tc := range testCases {
				rules := tc.rg.Rules(`Mozilla/5.0 (iPod touch; CPU iPhone OS 7_0_4 like Mac OS X)`, testDB.UAR.Regexes, tc.acc)
				if tc.exp == nil {
					Expect(rules).To(BeNil(), tc.desc)
				} else {
					Expect(rules).To(Equal(tc.exp), tc.desc)
				}
			}
		})

	})

	Describe("regexpSlice", func() {

		It("should unmarshal", func() {
			var s regexpSlice
			err := json.Unmarshal([]byte(`{"d":["A", "B", "C"], "6":{"1":"D"}}`), &s)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal(regexpSlice{
				regexp.MustCompile(`A`),
				regexp.MustCompile(`D`),
				regexp.MustCompile(`C`),
			}))
		})

	})

})

func iptr(i int) *int { return &i }
