package gildedrose

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/schema"
)

type SellByDateRangeJSON struct {
	UpperLimit *int `json:"upper_limit,omitempty"`
	LowerLimit *int `json:"lower_limit,omitempty"`
}

type GlobalPoliciesJSON struct {
	MaxQuality          int `json:"max_quality,omitempty"`
	MinQuality          int `json:"min_quality,omitempty"`
	ExpiredMultiplier   int `json:"expired_multiplier,omitempty"`
	DefaultQualityDelta int `json:"default_quality_delta,omitempty"`
}

type RulesJSON struct {
	ItemName         *string              `json:"item_name,omitempty"`
	RuleDescription  *string              `json:"rule_description,omitempty"`
	SellByDateRange  *SellByDateRangeJSON `json:"sell_by_date_range,omitempty"`
	Inmutable        *bool                `json:"inmutable,omitempty"`
	QualityDelta     *int                 `json:"quality_delta,omitempty"`
	ExpiresGradually *bool                `json:"expires_gradually,omitempty"`
}

type RulesTopLevelJSON struct {
	Schema         any                `json:"$schema,omitempty"`
	GlobalPolicies GlobalPoliciesJSON `json:"global_policies,omitempty"`
	Rules          []RulesJSON        `json:"rules,omitempty"`
}

func GenerateSchema() (string, error) {
	var dataType RulesTopLevelJSON
	s, err := schema.Generate(reflect.TypeOf(&dataType))
	if err != nil {
		return "", fmt.Errorf("failed to generate schema: %v", err)
	}
	jsonSchema, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshall schema to JSON: %v", err)
	}
	return string(jsonSchema), nil
}
