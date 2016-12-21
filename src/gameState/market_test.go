package gameState

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"util"
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

func testMarketTurn(t *testing.T, game *Game, playerName string, turn MarketTurn) []error {
	var backup *Game
	if iface, err := util.Copy(game); err != nil {
		panic(err)
	} else {
		backup = iface.(*Game)
	}
	if !reflect.DeepEqual(backup, game) {
		panic(fmt.Sprintf("fresh copy of the game is not equal\n\ncopy:%+v\n\noriginal:%+v\n",
			backup, game))
	}

	errs := game.PerformMarketTurn(playerName, turn)
	// If there were any errors at all the function should have done nothing to the game
	if len(errs) > 0 {
		if !reflect.DeepEqual(backup, game) {
			t.Errorf("game state changed from unsuccessful market turn\n\n%+v\n\n%+v",
				backup, game)
		}
	}
	return errs
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
	} else if errs := game.PerformMarketTurn("1st", MarketTurn{}); len(errs) == 0 {
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
	if errs := game.PerformMarketTurn("bad name", MarketTurn{}); len(errs) == 0 {
		t.Error("bad player name did not error performing market action")
	}

	for turn, actual := range game.TurnOrder {
		if turn != game.Turn {
			t.Fatalf("internal game turn %d != expected turn %d", game.Turn, turn)
		}
		for index := range rand.Perm(len(playerNames)) {
			if other := playerNames[index]; other != actual {
				if errs := game.PerformMarketTurn(other, MarketTurn{}); len(errs) == 0 {
					t.Errorf("%s's market action succeeded on %s's turn", other, actual)
				}
				if turn != game.Turn {
					t.Errorf("turn advanced from %s's action on %s's turn", other, actual)
				}
			}
		}
		if errs := game.PerformMarketTurn(actual, MarketTurn{}); len(errs) > 0 {
			t.Errorf("%s failed to perform their market action: %v", actual, errs)
		}
	}
}

// TestMarketActionValidation tests various MarketActions to make sure the ones lacking needed
// information and the ones whose information conflicts with the internal game state error.
func TestMarketActionValidation(t *testing.T) {
	game := NewGame([]string{"1st", "2nd", "3rd", "4th"})
	companyName := randomCompany()
	price := startingPrices[2][1]
	game.Companies[companyName].StockPrice = price

	// This just make the code shorter to fit the if statements on one line < 100 columns
	validateAction := func(action MarketAction) error {
		playerName := game.currentTurn()
		var errs []error
		if num := rand.Intn(2); num == 0 {
			errs = game.PerformMarketTurn(playerName, MarketTurn{Sales: []MarketAction{action}})
		} else {
			errs = game.PerformMarketTurn(playerName, MarketTurn{Purchase: &action})
		}
		if len(errs) == 0 {
			return nil
		} else if len(errs) > 1 {
			t.Errorf("market turn returned more errors than expected: %v", errs)
		}

		return errs[0]
	}

	if err := validateAction(MarketAction{Count: 1}); err == nil {
		t.Error("MarketAction with only count defined did not error")
	}
	if err := validateAction(MarketAction{Company: companyName}); err == nil {
		t.Error("MarketAction with only company defined did not error")
	}
	if err := validateAction(MarketAction{Company: "blah", Count: 1}); err == nil {
		t.Error("MarketAction with invalid company name did not error")
	}
	if err := validateAction(MarketAction{Company: companyName, Count: 1, Price: 3}); err == nil {
		t.Error("MarketAction with incorrect non-zero price did not error")
	}

	// Make sure the validation sets the price if we don't set it, and that the correct price
	// doesn't error when we provide it.
	action := MarketAction{Company: companyName, Count: 1}
	if _, err := game.validateAction(&action); err != nil {
		t.Errorf("market action %+v failed validation: %v", action, err)
	} else if action.Price != price {
		t.Error("action price $%d != expected stock price $%d", action.Price, price)
	} else if _, err = game.validateAction(&action); err != nil {
		t.Errorf("market action %+v failed validation: %v", action, err)
	}
}
