package devatlas

import (
	"encoding/json"
	"regexp"
	"strconv"
)

// UAR sub-tree
type UAR struct {
	Skip       []int        `json:"sk"`
	Regexes    regexpSlice  `json:"reg,omitempty"`
	RuleGroups []*RuleGroup `json:"rg"`
}

type RuleGroup struct {
	Sets []RuleSet `json:"t"`
	Reqs Indices   `json:"p"`
}

// Rules matches a list of rules for a given UA and pre-matched property-value indices
func (rg *RuleGroup) Rules(ua string, exps []*regexp.Regexp, acc Indices) []Rule {
	// Ensure that requirements are all met
	if len(rg.Reqs) == 0 {
		return nil
	}
	for pi, vi := range rg.Reqs {
		if acc[pi] != vi {
			return nil
		}
	}

	// If there is only one set, use its rules
	if len(rg.Sets) == 1 {
		return rg.Sets[0].Rules
	}

	// Merge matching rules
	rules := make([]Rule, 0, len(rg.Sets))
	for _, set := range rg.Sets {
		if set.Match(exps, ua) {
			rules = append(rules, set.Rules...)
		}
	}
	return rules
}

type RuleSet struct {
	Mi    *int   `json:"f"`
	Si    *int   `json:"s"`
	Rules []Rule `json:"r"`
}

func (rs *RuleSet) Match(exps []*regexp.Regexp, ua string) bool {
	if rs.Mi == nil {
		return false
	}
	mi := *rs.Mi
	return mi < len(exps) && exps[mi].MatchString(ua)
}

type Rule struct {
	Pi int  `json:"p"`
	Vi *int `json:"v"`
	Ri int  `json:"r"`
	M  int  `json:"m"`
}

type regexpSlice []*regexp.Regexp

func (ux *regexpSlice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Defaults  []string          `json:"d"`
		Overrides map[string]string `json:"6"` // set to API_VERSION
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	for _, expr := range raw.Defaults {
		x, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		*ux = append(*ux, x)
	}

	for s, expr := range raw.Overrides {
		x, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		i, _ := strconv.Atoi(s)
		(*ux)[i] = x
	}

	return nil
}
