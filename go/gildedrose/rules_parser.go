package gildedrose

import (
	"encoding/json"
	"fmt"
	"math"
)

var ZERO int = 0
var MINUS_ONE int = -1

func ParseRules(rawdata []byte) ([]Rule, error) {
	var data RulesTopLevelJSON
	err := json.Unmarshal(rawdata, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON data for Rules: %v", err)
	}

	var output []Rule

	for _, incomingRule := range data.Rules {

		qDelta := data.GlobalPolicies.DefaultQualityDelta
		if incomingRule.QualityDelta != nil {
			qDelta = *incomingRule.QualityDelta
		}

		bakedRule := NewSellByDateRule()
		var extraBakedRule Rule = nil

		if incomingRule.RuleDescription != nil {
			bakedRule.SetDescription(*incomingRule.RuleDescription)
		}

		// if the name is not explicitly defined in JSON rule, set it to be negative
		// when required, a new value for explicit negative check can be implemented easily
		if incomingRule.ItemName != nil {
			bakedRule.SetMatchName(*incomingRule.ItemName)
		} else {
			bakedRule.SetNegativeMatchName(incomingRule.ItemName == nil)
		}

		isInmutable := incomingRule.Inmutable != nil && *incomingRule.Inmutable
		if isInmutable {
			bakedRule.SetQualityDelta(0)
			bakedRule.SetSellInDelta(0)
			bakedRule.SetQualityUpperLimit(math.MaxInt)
			bakedRule.SetQualityLowerLimit(math.MinInt)
		} else {
			bakedRule.SetQualityDelta(qDelta)
			bakedRule.SetSellInDelta(-1)
			bakedRule.SetQualityUpperLimit(data.GlobalPolicies.MaxQuality)
			bakedRule.SetQualityLowerLimit(data.GlobalPolicies.MinQuality)
			// if a rule has some specific sell by rules, set them
			// if not, the global expire rule apply, by creating a duplicate
			if incomingRule.SellByDateRange != nil && !isInmutable {
				if incomingRule.SellByDateRange.LowerLimit != nil {
					bakedRule.SetSellByDateLowerLimit(*incomingRule.SellByDateRange.LowerLimit)
				}
				if incomingRule.SellByDateRange.UpperLimit != nil {
					bakedRule.SetSellByDateUpperLimit(*incomingRule.SellByDateRange.UpperLimit)
				}
			} else {
				expiredRule := bakedRule.Clone()
				expiredRule.SetDescription(fmt.Sprintf("[Expired] %s", expiredRule.Description))
				expiredRule.SetQualityDelta(
					data.GlobalPolicies.ExpiredMultiplier * qDelta)
				expiredRule.SetSellByDateUpperLimit(0)
				extraBakedRule = expiredRule
			}
		}

		if extraBakedRule != nil {
			output = append(output, extraBakedRule)
		}
		output = append(output, bakedRule)
		LOG.Printf("Parsed rules: \n%v\n%v", bakedRule, extraBakedRule)
	}

	return output, nil
}
