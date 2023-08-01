package gildedrose_test

import (
	"fmt"
	"testing"

	"github.com/emilybache/gildedrose-refactoring-kata/gildedrose"
)

type itemBeforeAfter struct {
	Name                        string
	BeforeSellIn, AfterSellIn   int
	BeforeQuality, AfterQuality int
}

func failTest[E any, G any](what string, item itemBeforeAfter, expected E, got G) error {
	return fmt.Errorf(
		"On item <<%v>>, %s does not match, expected %v but got %v",
		item,
		what,
		expected,
		got,
	)
}

func testUpdateQuality(testValues *[]itemBeforeAfter) error {
	var items = make([]*gildedrose.Item, len(*testValues))

	for i, testValue := range *testValues {
		items[i] = &gildedrose.Item{
			Name:    testValue.Name,
			Quality: testValue.BeforeQuality,
			SellIn:  testValue.BeforeSellIn,
		}
	}
	gildedrose.UpdateQuality(items)

	for i, testValue := range *testValues {
		if items[i].Name != testValue.Name {
			return failTest("Name", testValue, testValue.Name, items[i].Name)
		}
		if items[i].Quality != testValue.AfterQuality {
			return failTest("Quality", testValue, testValue.AfterQuality, items[i].Quality)
		}
		if items[i].SellIn != testValue.AfterSellIn {
			return failTest("SellIn", testValue, testValue.AfterSellIn, items[i].SellIn)
		}
	}
	return nil
}

func TestBasic(t *testing.T) {
	var basics = []itemBeforeAfter{
		{"+5 Dexterity Vest", 10, 9, 20, 19},
		{"Aged Brie", 2, 1, 0, 1},
		{"Elixir of the Mongoose", 5, 4, 7, 6},
		{"Sulfuras, Hand of Ragnaros", 0, 0, 80, 80},
		{"Sulfuras, Hand of Ragnaros", -1, -1, 80, 80},
		{"Backstage passes to a TAFKAL80ETC concert", 15, 14, 20, 21},
		{"Backstage passes to a TAFKAL80ETC concert", 10, 9, 49, 50},
		{"Backstage passes to a TAFKAL80ETC concert", 5, 4, 49, 50},
	}

	err := testUpdateQuality(&basics)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestNotReallyConjured(t *testing.T) {
	var conjuredItems = []itemBeforeAfter{
		{"Conjured Not Really Cake", 10, 9, 6, 5},
	}
	err := testUpdateQuality(&conjuredItems)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestConjured(t *testing.T) {
	var conjuredItems = []itemBeforeAfter{
		{"Conjured Mana Cake", 3, 2, 6, 4},     // normal conjured
		{"Conjured Mana Cake", 50, 49, 50, 48}, // conjured at 50 quality
		{"Conjured Mana Cake", 10, 9, 1, -1},   // go to negative quality
		{"Conjured Mana Cake", -2, -3, 10, 6},  // expired
	}

	err := testUpdateQuality(&conjuredItems)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestBackstagePasses(t *testing.T) {
	var pass = []itemBeforeAfter{
		{"Backstage passes to a TAFKAL80ETC concert", 9, 8, 10, 12},
		{"Backstage passes to a TAFKAL80ETC concert", 2, 1, 10, 13},
		{"Backstage passes to a TAFKAL80ETC concert", 0, -1, 50, 0},
	}
	err := testUpdateQuality(&pass)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestExpired(t *testing.T) {
	var item = []itemBeforeAfter{
		{"Chorizo", -10, -11, 10, 8},
		{"Aged Brie", -1, -2, 20, 22},
	}
	err := testUpdateQuality(&item)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}
