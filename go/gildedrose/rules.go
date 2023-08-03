package gildedrose

import (
	"math"
	"strings"
)

type Rule interface {
	Apply(*Item) bool
}

// SellByDateRule describes a change in a item
// by matching its name, and that the SellByDateLimit is lower
// than the current value

func NewSellByDateRule() *SellByDateRule {
	return &SellByDateRule{
		NegativeMatchName: false,
	}
}

type SellByDateRule struct {
	Description          string
	MatchName            *string
	NegativeMatchName    bool
	SellByDateUpperLimit *int
	SellByDateLowerLimit *int
	QualityDelta         *int
	SellInDelta          *int
	QualityUpperLimit    *int
	QualityLowerLimit    *int
}

// golang doesn't have this on stdlib???
func clone[T any](in *T) *T {
	if in == nil {
		return nil
	}

	var new T
	new = *in
	return &new
}

func (r *SellByDateRule) Clone() *SellByDateRule {
	var new SellByDateRule
	new.Description = r.Description
	new.MatchName = clone(r.MatchName)
	new.NegativeMatchName = r.NegativeMatchName
	new.QualityDelta = clone(r.QualityDelta)
	new.SellInDelta = clone(r.SellInDelta)
	new.QualityLowerLimit = clone(r.QualityLowerLimit)
	new.QualityUpperLimit = clone(r.QualityUpperLimit)
	new.SellByDateLowerLimit = clone(r.SellByDateLowerLimit)
	new.SellByDateUpperLimit = clone(r.SellByDateUpperLimit)
	return &new
}

func (r *SellByDateRule) SetDescription(desc string) *SellByDateRule {
	r.Description = desc
	return r
}

func (r *SellByDateRule) SetMatchName(name string) *SellByDateRule {
	r.MatchName = &name
	return r
}

func (r *SellByDateRule) SetNegativeMatchName(enabled bool) *SellByDateRule {
	r.NegativeMatchName = enabled
	return r
}

func (r *SellByDateRule) SetSellByDateLowerLimit(lower int) *SellByDateRule {
	r.SellByDateLowerLimit = &lower
	return r
}

func (r *SellByDateRule) SetSellByDateUpperLimit(upper int) *SellByDateRule {
	r.SellByDateUpperLimit = &upper
	return r
}

func (r *SellByDateRule) SetQualityUpperLimit(upper int) *SellByDateRule {
	r.QualityUpperLimit = &upper
	return r
}

func (r *SellByDateRule) SetQualityLowerLimit(lower int) *SellByDateRule {
	r.QualityLowerLimit = &lower
	return r
}

func (r *SellByDateRule) SetSellByDateLimits(upper, lower int) *SellByDateRule {
	return r.SetSellByDateLowerLimit(lower).SetSellByDateUpperLimit(upper)
}

func (r *SellByDateRule) SetQualityDelta(delta int) *SellByDateRule {
	r.QualityDelta = &delta
	return r
}

func (r *SellByDateRule) SetSellInDelta(delta int) *SellByDateRule {
	r.SellInDelta = &delta
	return r
}

func (r *SellByDateRule) nameMatches(item *Item) bool {
	if r.MatchName == nil {
		return true
	}
	gotMatch := strings.EqualFold(*r.MatchName, item.Name)
	return gotMatch != r.NegativeMatchName
}

func (r *SellByDateRule) isInSellByDateRange(item *Item) bool {
	var lower int
	var upper int

	if r.SellByDateLowerLimit == nil {
		lower = math.MinInt
	} else {
		lower = *r.SellByDateLowerLimit
	}

	if r.SellByDateUpperLimit == nil {
		upper = math.MaxInt
	} else {
		upper = *r.SellByDateUpperLimit
	}

	return (item.SellIn <= upper) && (item.SellIn >= lower)
}

func (r *SellByDateRule) matches(item *Item) bool {
	return (r.nameMatches(item) && r.isInSellByDateRange(item))
}

func (r *SellByDateRule) enforceQualityLimits(item *Item) {
	if r.QualityLowerLimit != nil && *r.QualityLowerLimit > item.Quality {
		item.Quality = *r.QualityLowerLimit
	}
	if r.QualityUpperLimit != nil && *r.QualityUpperLimit < item.Quality {
		item.Quality = *r.QualityUpperLimit
	}
}

func (r *SellByDateRule) calculateNewValues(item *Item) {
	if r.QualityDelta != nil {
		item.Quality += *r.QualityDelta
	}
	if r.SellInDelta != nil {
		item.SellIn += *r.SellInDelta
	}
}

func (r *SellByDateRule) Apply(item *Item) bool {
	if !r.matches(item) {
		return false
	}
	LOG.Printf("Rule <<%s>> matches %v", r.Description, item)

	r.calculateNewValues(item)
	r.enforceQualityLimits(item)

	return true
}

func ApplySet(ruleSet []Rule, item *Item) {
	for _, r := range ruleSet {
		matched := r.Apply(item)
		if matched {
			break
		}
	}
}
