package devatlas

// UAR extensions
type UAR struct {
	Skip       []int        `json:"sk"`
	Regexp     regexpSliceO `json:"reg,omitempty"`
	RuleGroups []ruleGroup  `json:"rg"`
}

func (u *UAR) update(ua string, acc indexMap) map[int]interface{} {
	// Check if any of the pre-resolved properties require skipping UAR traversal
	for _, pi := range u.Skip {
		if acc[pi] > 0 {
			return nil
		}
	}

	res := make(map[int]interface{}, len(u.RuleGroups))
	for _, rg := range u.RuleGroups {
		for _, rule := range rg.matchRules(ua, u.Regexp, acc) {
			if rule.Vi != nil {
				acc[rule.Pi] = *rule.Vi
			} else if m := u.Regexp[rule.Ri].FindStringSubmatch(ua); rule.M < len(m) {
				res[rule.Pi] = m[rule.M]
			}
		}
	}
	return res
}

// Plain rule
type rule struct {
	Pi int  `json:"p"`
	Vi *int `json:"v"`
	Ri int  `json:"r"`
	M  int  `json:"m"`
}

// A rule-set combines multiple rules
type ruleSet struct {
	Mi    *int   `json:"f"`
	Si    *int   `json:"s"`
	Rules []rule `json:"r"`
}

func (rs *ruleSet) Match(exps []*Regexp, ua string) bool {
	if rs.Mi == nil {
		return false
	}
	mi := *rs.Mi
	return mi < len(exps) && exps[mi].MatchString(ua)
}

// A rule-group combines multiple rule-sets

type ruleGroup struct {
	Sets []ruleSet `json:"t"`
	Reqs indexMap  `json:"p"`
}

// matches a list of rules for a given UA and pre-matched property-value indices
func (rg *ruleGroup) matchRules(ua string, exps []*Regexp, acc indexMap) []rule {
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
	rules := make([]rule, 0, len(rg.Sets))
	for _, set := range rg.Sets {
		if set.Match(exps, ua) {
			rules = append(rules, set.Rules...)
		}
	}
	return rules
}
