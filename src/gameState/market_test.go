package gameState

import (
	"math/rand"
	"testing"
)

// TestMarketActionNonMarketPhase checks to make sure market actions can only be performed during
// the market phase.
func TestMarketActionPhaseValidation(t *testing.T) {
	game := NewGame([]string{"1st", "2nd", "3rd", "4th"})
	game.beginBusinessPhase()

	if game.marketPhase() {
		t.Fatal("game unexpected in market phase")
	} else if !game.businessPhase() {
		t.Fatal("game not in a valid phase state")
	} else if err := game.PerformMarketAction("blah", MarketAction{}); err == nil {
		t.Error("market action did not fail while in the business phase")
	}
}

// TestMarketPlayerValidation checks the error when a non-existent players tries to perform a market
// action. It also checks to make sure no player can perform a market action unless its their turn.
func TestMarketPlayerValidation(t *testing.T) {
	playerNames := []string{"1st", "2nd", "3rd", "4th", "5th", "6th"}
	game := NewGame(playerNames)

	// This part of the test must come first, otherwise we will no longer be in the market phase
	// and this would no longer check player name existence.
	if err := game.PerformMarketAction("bad name", MarketAction{}); err == nil {
		t.Error("bad player name did not error performing market action")
	}

	for turn, actual := range game.turnOrder {
		if turn != game.turn {
			t.Fatalf("internal game turn %d != expected turn %d", game.turn, turn)
		}
		for index := range rand.Perm(len(playerNames)) {
			if other := playerNames[index]; other != actual {
				if err := game.PerformMarketAction(other, MarketAction{}); err == nil {
					t.Errorf("%s's market action succeeded on %s's turn", other, actual)
				}
				if turn != game.turn {
					t.Errorf("turn advanced from %s's action on %s's turn", other, actual)
				}
			}
		}
		if err := game.PerformMarketAction(actual, MarketAction{}); err != nil {
			t.Errorf("%s failed to perform their market action: %v", actual, err)
		}
	}
}

// TestMarketActionValidation tests various MarketActions to make sure the ones lacking needed
// information and the ones whose information conflicts with the internal game state error.
func TestMarketActionValidation(t *testing.T) {
	game := NewGame([]string{"1st", "2nd", "3rd", "4th"})
	var companyName string
	for name, _ := range game.Companies {
		companyName = name
	}
	game.Companies[companyName].StockPrice = 100
	playerName := game.turnOrder[0]

	// This just make the code shorter to fit the if statements on one line < 100 columns
	performAction := func(action MarketAction) error {
		return game.PerformMarketAction(playerName, action)
	}

	if err := performAction(MarketAction{Count: 1}); err == nil {
		t.Error("MarketAction with only count defined did not error")
	}
	if err := performAction(MarketAction{Company: companyName}); err == nil {
		t.Error("MarketAction with only company defined did not error")
	}
	if err := performAction(MarketAction{Company: "blah", Count: 1}); err == nil {
		t.Error("MarketAction with invalid company name did not error")
	}
	if err := performAction(MarketAction{Company: companyName, Count: 1, Price: 3}); err == nil {
		t.Error("MarketAction with incorrect non-zero price did not error")
	}

	// Lastly make sure that we can perform a valid action to make sure we weren't accidentally
	// failing earlier due to something we weren't expecting.
	if err := performAction(MarketAction{Company: companyName, Count: 1, Price: 100}); err != nil {
		t.Errorf("%s failed to buy 1 share in %s: %v", playerName, companyName, err)
	}
}

// startCompany is a convenience function for the tests that check starting companies under various
// conditions. It creates a new game for every call, initializes some of the specified parameters,
// then performs the market action with the remaining parameters.
func startCompany(companyName string, count, price, cash, techLvl int) (error, int) {
	game := NewGame([]string{"player"})
	game.TechLevel = techLvl
	game.Players["player"].Cash = cash
	err := game.PerformMarketAction("player", MarketAction{
		Company: companyName,
		Count:   count,
		Price:   price,
	})
	return err, game.Companies[companyName].StockPrice
}

// startingPrices is the list of valid stock prices for each tech level (as printed on the physical
// board of the original game).
var startingPrices = [][3]int{
	[3]int{55, 60, 66},
	[3]int{60, 66, 74},
	[3]int{66, 74, 82},
	[3]int{74, 82, 91},
	[3]int{82, 91, 100},
}

// TestStartingTech3Company makes sure the companies labeled at only available after tech level 3
// cannot be started before that.
func TestStartingTech3Company(t *testing.T) {
	var companyName string
	for name, start := range companyInitCond {
		if start.tech3 {
			companyName = name
			break
		}
	}
	for ind, prices := range startingPrices[:2] {
		techLvl := ind + 1
		price := prices[rand.Intn(3)]
		if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err == nil {
			t.Errorf("starting %s during tech level %d did not error", companyName, techLvl)
		} else if endPrice != 0 {
			t.Errorf("failed starting %s in tech level %d left price as %d$",
				companyName, techLvl, endPrice)
		}
	}
	for ind, prices := range startingPrices[2:] {
		techLvl := ind + 3
		price := prices[rand.Intn(3)]
		if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err != nil {
			t.Errorf("starting %s during tech level %d failed: %v", companyName, techLvl, err)
		} else if endPrice != price {
			t.Errorf("starting %s in tech level %d at %d$ left price as %d$",
				companyName, techLvl, price, endPrice)
		}
	}
}

// TestStartingCompanyPrices checks to make sure a company can be started at any price valid for
// the current tech level, and none other.
func TestStartingCompanyPrices(t *testing.T) {
	var companyName string
	for name, start := range companyInitCond {
		if !start.tech3 {
			companyName = name
			break
		}
	}

	for ind, prices := range startingPrices {
		techLvl := ind + 1
		for _, price := range prices {
			if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err != nil {
				t.Errorf("starting %s in tech level %d at %d$ failed: %v",
					companyName, techLvl, price, err)
			} else if endPrice != price {
				t.Errorf("starting %s in tech level %d at %d$ left price as %d$",
					companyName, techLvl, price, endPrice)
			}
		}

		invalidPrices := []int{
			prices[0] - (rand.Intn(20) + 1),
			prices[0] + rand.Intn(prices[1]-prices[0]-1) + 1,
			prices[1] + rand.Intn(prices[2]-prices[1]-1) + 1,
			prices[2] + rand.Intn(20) + 1,
		}
		for _, price := range invalidPrices {
			if err, endPrice := startCompany(companyName, 1, price, price*2, techLvl); err == nil {
				t.Errorf("starting %s in tech level %d at %d$ (invalid) did not error",
					companyName, techLvl, price)
			} else if endPrice != 0 {
				t.Errorf("failed starting %s in tech level %d at %d$ (invalid) left price as %d$",
					companyName, techLvl, price, endPrice)
			}
		}
	}
}

// TestStartingCompanyFailure checks to make sure a company's stock price doesn't change if the
// player doesn't have enough cash to complete the transaction.
func TestStartingCompanyFailure(t *testing.T) {
	var companyName string
	for name, start := range companyInitCond {
		if !start.tech3 {
			companyName = name
			break
		}
	}
	price := startingPrices[0][1]
	count := rand.Intn(9) + 1
	if err, endPrice := startCompany(companyName, count, price, count*price-20, 1); err == nil {
		t.Errorf("attempt to buy %d %s shares at %d$ with insufficient cash did not error",
			count, companyName, price)
	} else if endPrice != 0 {
		t.Errorf("attempt to buy %d %s shares at %d$ with insufficient cash left price at %d$",
			count, companyName, price, endPrice)
	}
}
