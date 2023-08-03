package gildedrose

import "log"

type Item struct {
	Name            string
	SellIn, Quality int
}

var LOG *log.Logger = log.Default()

var loadedRules []Rule

func SetCurrentRules(ruleSet []Rule) {
	loadedRules = ruleSet
}

func UpdateQuality(items []*Item) {
	// if no rules are loaded, bail out
	if loadedRules == nil {
		LOG.Fatalf("No rules have been installed!")
		return
	}
	for _, item := range items {
		ApplySet(loadedRules, item)
	}
}
