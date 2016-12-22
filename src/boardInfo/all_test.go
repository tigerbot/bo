package boardInfo_test

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	. "boardInfo"
)

func TestBuildCost(t *testing.T) {
	type testPair struct {
		hexcoord string
		cost     int
	}
	testVals := []testPair{
		{hexcoord: "A30", cost: 40},
		{hexcoord: "E24", cost: 40},
		{hexcoord: "G20", cost: 20},
		{hexcoord: "H1", cost: 10},
		{hexcoord: "I16", cost: 100},
		{hexcoord: "J17", cost: 80},
		{hexcoord: "A26", cost: 0},
		{hexcoord: "J16", cost: 0},
		{hexcoord: "J18", cost: 0},
	}

	for _, pair := range testVals {
		if cost := BuildCost(pair.hexcoord); cost != pair.cost {
			t.Errorf("expected hex %q to cost $%d, but got $%d", pair.hexcoord, pair.cost, cost)
		}
	}
}

func TestTrainCost(t *testing.T) {
	type testPair struct{ trainNum, cost int }
	testVals := []testPair{
		{trainNum: 1, cost: 100},
		{trainNum: 5, cost: 80},
		{trainNum: 6, cost: 140},
		{trainNum: 13, cost: 170},
		{trainNum: 17, cost: 260},
		{trainNum: 21, cost: 380},
		{trainNum: 26, cost: 500},
		{trainNum: 30, cost: 380},
	}

	for _, pair := range testVals {
		if cost := TrainCost(pair.trainNum); cost != pair.cost {
			t.Errorf("expected train #%d to cost $%d, but got $%d", pair.trainNum, pair.cost, cost)
		}
	}
}

func TestTechLevel(t *testing.T) {
	type testPair struct{ trainsSold, level int }
	testVals := []testPair{
		{trainsSold: 0, level: 1},
		{trainsSold: 5, level: 1},
		{trainsSold: 6, level: 2},
		{trainsSold: 8, level: 2},
		{trainsSold: 10, level: 2},
		{trainsSold: 11, level: 3},
		{trainsSold: 14, level: 3},
		{trainsSold: 16, level: 4},
		{trainsSold: 25, level: 5},
		{trainsSold: 26, level: 6},
	}

	for _, pair := range testVals {
		if level := TechLevel(pair.trainsSold); level != pair.level {
			t.Errorf("expected tech level %d after %d trains sold, but got %d",
				pair.level, pair.trainsSold, level)
		}
	}
}

func TestTileAdjacency(t *testing.T) {
	type testSet struct {
		coordA, coordB string
		adjacent       bool
	}
	testVals := []testSet{
		// hexes are not adjacent to themselves.
		{adjacent: false, coordA: "A30", coordB: "A30"},
		// error conditions
		{adjacent: false, coordA: "IJK", coordB: "H24"},
		{adjacent: false, coordA: "F9", coordB: "GONE"},
		// exceptions from the normal rules
		{adjacent: false, coordA: "I22", coordB: "I24"},
		{adjacent: false, coordA: "I24", coordB: "I22"},
		{adjacent: false, coordA: "E12", coordB: "F13"},
		{adjacent: false, coordA: "F13", coordB: "E12"},
		// check adjacency in all six directions
		{adjacent: true, coordA: "B27", coordB: "B29"}, // right
		{adjacent: true, coordA: "I10", coordB: "J11"}, // down-right
		{adjacent: true, coordA: "A30", coordB: "B29"}, // down-left
		{adjacent: true, coordA: "G6", coordB: "G4"},   // left
		{adjacent: true, coordA: "C28", coordB: "B27"}, // up-left
		{adjacent: true, coordA: "H9", coordB: "G10"},  // up-right

		{adjacent: false, coordA: "C26", coordB: "B29"},
		{adjacent: false, coordA: "C26", coordB: "E26"},
	}

	for _, set := range testVals {
		if adjacent := TilesAdjacent(set.coordA, set.coordB); adjacent != set.adjacent {
			if set.adjacent {
				t.Errorf("expected %q and %q to be adjacent", set.coordA, set.coordB)
			} else {
				t.Errorf("expected %q and %q to not be adjacent", set.coordA, set.coordB)
			}
		}
	}
}

func TestStartingCity(t *testing.T) {
	type testPair struct{ company, hexcoord string }
	testVals := []testPair{
		{company: "Erie", hexcoord: "D19"},
		{company: "New York Central", hexcoord: "D25"},
		{company: "New York, Chicago & Saint Louis", hexcoord: "E4"},
		{company: "Pennsylvania", hexcoord: "G24"},
		{company: "Boston & Maine", hexcoord: "D29"},
		{company: "New York, New Haven & Hartford", hexcoord: "E28"},
		{company: "Wabash", hexcoord: "G8"},
		{company: "Baltimore & Ohio", hexcoord: "H23"},
		{company: "Illinois Central", hexcoord: "I0"},
		{company: "Chesapeake & Ohio", hexcoord: "J21"},
	}

	for _, pair := range testVals {
		if value := StartingLocation(pair.company); value != pair.hexcoord {
			t.Errorf("expected %s to start on %q, not %q", pair.company, pair.hexcoord, value)
		}
	}
}

func TestCityList(t *testing.T) {
	type testPair struct{ hexes, cities []string }
	testVals := []testPair{
		{
			hexes:  []string{"C22", "D21", "E20", "F19", "G18", "H17"},
			cities: []string{"Syracuse"},
		},
		{
			hexes:  []string{"G24", "G22", "G20", "G18", "G16", "G14"},
			cities: []string{"Pittsburgh", "Harrisburg", "Philadelphia"},
		},
	}

	for _, pair := range testVals {
		cityList := Cities(pair.hexes...)
		cityNames := make([]string, 0, len(cityList))
		for _, city := range cityList {
			cityNames = append(cityNames, city.Name)
		}
		sort.Strings(cityNames)
		sort.Strings(pair.cities)
		if !reflect.DeepEqual(cityNames, pair.cities) {
			t.Errorf("expected cities %q in hexes %q, got %q", pair.cities, pair.hexes, cityNames)
		}
	}
}

func TestStockIncrease(t *testing.T) {
	if value := NextStockPrice(375); value != 375 {
		t.Errorf("NextStockPrice did not stall at max value 375: %d", value)
	}
	if value := NextStockPrice(500); value != 375 {
		t.Errorf("NextStockPrice did not bring bad price 500 to max value 375: %d", value)
	}
	prices := StartingStockPrices(rand.Intn(5) + 1)
	if value := NextStockPrice(prices[1]); value != prices[2] {
		t.Errorf("expected stock price after %d to be %d, got %d", prices[1], prices[2], value)
	}
	badPrice := prices[0] + rand.Intn(prices[1]-prices[0]-1) + 1
	if value := NextStockPrice(badPrice); value != prices[1] {
		t.Errorf("expected stock price after bad value %d to be %d, got %d",
			badPrice, prices[1], value)
	}
}

func TestStockDecreate(t *testing.T) {
	if value := PrevStockPrice(34); value != 34 {
		t.Errorf("PrevStockPrice did not stall at min value 34: %d", value)
	}
	if value := PrevStockPrice(0); value != 34 {
		t.Errorf("PrevStockPrice did not bring bad price 0 to min value 34: %d", value)
	}
	prices := StartingStockPrices(rand.Intn(5) + 1)
	if value := PrevStockPrice(prices[1]); value != prices[0] {
		t.Errorf("expected stock price before %d to be %d, got %d", prices[1], prices[0], value)
	}
	badPrice := prices[0] + rand.Intn(prices[1]-prices[0]-1) + 1
	if value := PrevStockPrice(badPrice); value != prices[0] {
		t.Errorf("expected stock price before bad value %d to be %d, got %d",
			badPrice, prices[0], value)
	}
}
