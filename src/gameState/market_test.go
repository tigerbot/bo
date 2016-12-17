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

// TestCompanyStartTech3 makes sure the companies labeled at only available after tech level 3
// cannot be started before that.
func TestCompanyStartTech3(t *testing.T) {
	companyName := randomCompany(true)

	for ind, prices := range startingPrices[:2] {
		techLvl := ind + 1
		price := prices[rand.Intn(3)]
		if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err == nil {
			t.Errorf("starting %s during tech level %d did not error", companyName, techLvl)
		} else if endPrice != 0 {
			t.Errorf("failed starting %s in tech level %d left price as $%d",
				companyName, techLvl, endPrice)
		}
	}

	for ind, prices := range startingPrices[2:] {
		techLvl := ind + 3
		price := prices[rand.Intn(3)]
		if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err != nil {
			t.Errorf("starting %s during tech level %d failed: %v", companyName, techLvl, err)
		} else if endPrice != price {
			t.Errorf("starting %s in tech level %d at $%d left price as $%d",
				companyName, techLvl, price, endPrice)
		}
	}
}

// TestCompanyStartPrices checks to make sure a company can be started at any price valid for
// the current tech level, and none other.
func TestCompanyStartPrices(t *testing.T) {
	companyName := randomCompany(false)

	for ind, prices := range startingPrices {
		techLvl := ind + 1

		for _, price := range prices {
			if err, endPrice := startCompany(companyName, 4, price, 500, techLvl); err != nil {
				t.Errorf("starting %s in tech level %d at $%d failed: %v",
					companyName, techLvl, price, err)
			} else if endPrice != price {
				t.Errorf("starting %s in tech level %d at $%d left price as $%d",
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
				t.Errorf("starting %s in tech level %d at $%d (invalid) did not error",
					companyName, techLvl, price)
			} else if endPrice != 0 {
				t.Errorf("failed starting %s in tech level %d at $%d (invalid) left price as $%d",
					companyName, techLvl, price, endPrice)
			}
		}
	}
}

// TestCompanyStartFailure checks to make sure a company's stock price doesn't change if the
// player doesn't have enough cash to complete the transaction.
func TestCompanyStartFailure(t *testing.T) {
	companyName := randomCompany(false)
	price := startingPrices[0][1]
	count := rand.Intn(9) + 1

	if err, endPrice := startCompany(companyName, count, price, count*price-20, 1); err == nil {
		t.Errorf("attempt to buy %d %s shares at $%d with insufficient cash did not error",
			count, companyName, price)
	} else if endPrice != 0 {
		t.Errorf("attempt to buy %d %s shares at $%d with insufficient cash left price at $%d",
			count, companyName, price, endPrice)
	}
}

// TestStockBuyOrphan checks to make sure that if a company has orphaned stock when a player tries
// to buy its stock, the orphan stock is bought before the company's held stock.
func TestStockBuyOrphan(t *testing.T) {
	orphaned := rand.Intn(7) + 2
	held := rand.Intn(10 - orphaned)
	stockPrice := 60
	extraCash := rand.Intn(200)

	companyName := randomCompany()
	game := NewGame([]string{"player"})
	buyStock := func(count int) error {
		game.OrphanStocks[companyName] = orphaned
		game.Companies[companyName].HeldStock = held
		game.Companies[companyName].Treasury = 0
		game.Companies[companyName].StockPrice = stockPrice

		startingCash := stockPrice*count + extraCash
		game.Players["player"].Cash = startingCash
		game.Players["player"].Stocks[companyName] = 0

		// If the transaction succeeded, the effect on the player should be the same no matter
		// where the stock they purchased came from. If the transaction failed there should have
		// been no affect on the player.
		err := game.PerformMarketAction("player", MarketAction{Company: companyName, Count: count})
		if err == nil {
			if playerHeld := game.Players["player"].Stocks[companyName]; playerHeld != count {
				t.Errorf("player held stock = %d after buying %d stocks", playerHeld, count)
			}
			if playerCash := game.Players["player"].Cash; playerCash != extraCash {
				t.Errorf("player cash $%d after transaction, expected $%d", playerCash, extraCash)
			}
		} else {
			if playerHeld := game.Players["player"].Stocks[companyName]; playerHeld != 0 {
				t.Errorf("failed purchase of %d stock gain player %d stock", count, playerHeld)
			}
			if playerCash := game.Players["player"].Cash; playerCash != startingCash {
				t.Errorf("failed purchase of %d stock at $%d cost player $%d",
					count, stockPrice, startingCash-playerCash)
			}
		}
		return err
	}

	count := orphaned + held + rand.Intn(10-(orphaned+held)) + 1
	if err := buyStock(count); err == nil {
		t.Error("attempt to buy more stock (%d) than available (%d+%d) did not error",
			count, held, orphaned)
	}

	count = orphaned + rand.Intn(held+1)
	if err := buyStock(count); err != nil {
		t.Errorf("attempt to buy all orphaned stock (%d >= %d) failed: %v", count, orphaned, err)
	} else {
		if left, present := game.OrphanStocks[companyName]; left != 0 {
			t.Errorf("attempt to buy all orphaned stock (%d >= %d) left %d orphaned stock",
				count, orphaned, left)
		} else if present {
			t.Error("buying all orphaned stock did not remove company name from the map")
		}
		if left := game.Companies[companyName].HeldStock; left != orphaned+held-count {
			t.Errorf("remaining held stock %d != expected %d+%d - %d", left, orphaned, held, count)
		}
		if game.Companies[companyName].Treasury != stockPrice*(count-orphaned) {
			t.Errorf("company gained $%d from purchase of %d stock at $%d with %d orphaned",
				game.Companies[companyName].Treasury, count, stockPrice, orphaned)
		}
	}

	count = rand.Intn(orphaned-1) + 1
	if err := buyStock(count); err != nil {
		t.Errorf("attempt to buy some orphaned stock (%d) failed: %v", count, err)
	} else {
		if left := game.OrphanStocks[companyName]; left != orphaned-count {
			t.Errorf("remaining orphaned stock %d != expected %d - %d", left, orphaned, count)
		}
		if left := game.Companies[companyName].HeldStock; left != held {
			t.Errorf("purchase of %d/%d orphaned stock changed company held stock (%d->%d)",
				count, orphaned, held, left)
		}
		if game.Companies[companyName].Treasury != 0 {
			t.Errorf("purchase of %d/%d orphaned stock gave company $%d",
				count, orphaned, game.Companies[companyName].Treasury)
		}
	}
}
