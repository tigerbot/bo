package gameState

import (
	"math/rand"
	"testing"
)

func randomCompany(techLvl3 ...bool) string {
	// We have no control over the order we iterate over a map, but it's not necessarily random.
	// So in order to do true random we create a list of company names that match the requirement
	// and then return a random value from that list
	nameList := make([]string, 0, len(companyInitCond))
	for name, start := range companyInitCond {
		if len(techLvl3) == 0 {
			nameList = append(nameList, name)
		} else if techLvl3[0] == start.tech3 {
			nameList = append(nameList, name)
		}
	}
	return nameList[rand.Intn(len(nameList))]
}

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
	companyName := randomCompany()
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
