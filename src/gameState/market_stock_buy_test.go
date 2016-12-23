package gameState

import (
	"math/rand"
	"testing"
)

func buyStock(t *testing.T, game *Game, playerName, companyName string, count int) []error {
	startPrice := game.Companies[companyName].StockPrice
	startHeld := game.Companies[companyName].HeldStock
	startTreasure := game.Companies[companyName].Treasury
	startOrphan := game.OrphanStocks[companyName]

	startStock := game.Players[playerName].Stocks[companyName]
	startCash := game.Players[playerName].Cash
	totalCost := count * startPrice

	action := MarketAction{Company: companyName, Count: count}
	if errs := testMarketTurn(t, game, playerName, MarketTurn{Purchase: &action}); len(errs) > 0 {
		return errs
	}

	if stockPrice := game.Companies[companyName].StockPrice; stockPrice != startPrice {
		t.Errorf("buying stock changed price from $%d to $%d", startPrice, stockPrice)
	}
	if orphan := game.OrphanStocks[companyName]; orphan != startOrphan {
		t.Errorf("buying stock changed changed orphan stock from %d to %d", startOrphan, orphan)
	}
	if change := game.Players[playerName].Stocks[companyName] - startStock; change != count {
		t.Errorf("player gained %d stock after buying %d", change, count)
	}
	if change := startCash - game.Players[playerName].Cash; change != totalCost {
		t.Errorf("player lost $%d after buying %d stock at $%d", change, count, startPrice)
	}
	if change := startHeld - game.Companies[companyName].HeldStock; change != count {
		t.Errorf("company lost %d stock after player bought %d", change, count)
	}
	if change := game.Companies[companyName].Treasury - startTreasure; change != totalCost {
		t.Errorf("company gained $%d after player bought %d stock at $%d",
			change, count, startPrice)
	}
	return nil
}

// TestStockBuy checks to make sure players cannot buy more stock than the they can afford or
// more stock than the company has, and that the stock and money changes appropriately when the
// purchase is valid.
func TestStockBuy(t *testing.T) {
	stockPrice := startingPrices[0][1]
	companyName := randomCompany()
	testPurchase := func(count, startingHeld, startingCash int) []error {
		game := NewGame([]string{"player"})
		game.Companies[companyName].StockPrice = stockPrice
		game.Companies[companyName].HeldStock = startingHeld
		game.Companies[companyName].Treasury = 0

		game.Players["player"].Cash = startingCash
		game.Players["player"].Stocks[companyName] = 0
		return buyStock(t, game, "player", companyName, count)
	}

	if errs := testPurchase(6, 4, 600); len(errs) == 0 {
		t.Error("attempt to buy more stock than the company held succeeded")
	}
	if errs := testPurchase(6, 8, 300); len(errs) == 0 {
		t.Error("attempt to buy more stock than player can afford succeeded")
	}
	if errs := testPurchase(4, 6, 250); len(errs) > 0 {
		t.Errorf("attempt to buy stock failed unexpectedly: %v", errs)
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
	// We require a separate test function because of the affect of the orphan stock on where stock
	// comes from and where the money spent of them goes.
	buyOrphanStock := func(count int) []error {
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
		action := MarketAction{Company: companyName, Count: count}
		errs := game.PerformMarketTurn("player", MarketTurn{Purchase: &action})
		if len(errs) == 0 {
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
		return errs
	}

	count := orphaned + held + rand.Intn(10-(orphaned+held)) + 1
	if errs := buyOrphanStock(count); len(errs) == 0 {
		t.Error("attempt to buy more stock (%d) than available (%d+%d) did not error",
			count, held, orphaned)
	}

	count = orphaned + rand.Intn(held+1)
	if errs := buyOrphanStock(count); len(errs) > 0 {
		t.Errorf("attempt to buy all orphaned stock (%d >= %d) failed: %v", count, orphaned, errs)
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
	if errs := buyOrphanStock(count); len(errs) > 0 {
		t.Errorf("attempt to buy some orphaned stock (%d) failed: %v", count, errs)
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

// TestStockBuyTakeover checks the transitional between presidents when a player acquires more
// stock in a company than the previous president.
func TestStockBuyTakeover(t *testing.T) {
	// Create a new game and override the random turn order to make the test easier
	game := NewGame([]string{"pres", "other"})

	action := MarketAction{
		Company: randomCompany(false),
		Price:   startingPrices[0][1],
		Count:   4,
	}
	// Make sure the players have enough cash to purchase all the stock needed
	game.Players["pres"].Cash = 4 * action.Price
	game.Players["other"].Cash = 6 * action.Price

	// Establish the initial president the other player will try to take over. Turn order changes
	// are just to make it so we don't need the president to pass later.
	game.TurnManager.Order = []string{"pres"}
	if errs := startCompany(t, game, action.Company, action.Count, action.Price); len(errs) > 0 {
		t.Fatalf("failed to establish initial president: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("expected company president to be pres, but is instead %q", president)
	}

	game.TurnManager.Order = []string{"other"}
	if errs := buyStock(t, game, "other", action.Company, 2); len(errs) > 0 {
		t.Fatalf("other player failed to purchase additional stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("president changed to %q when other player has less stock", president)
	}

	if errs := buyStock(t, game, "other", action.Company, 2); len(errs) > 0 {
		t.Fatalf("other player failed to purchase additional stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("president changed to %q when other player tied for stock", president)
	}

	if errs := buyStock(t, game, "other", action.Company, 2); len(errs) > 0 {
		t.Fatalf("other player failed to purchase additional stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "other" {
		t.Fatalf("president changed to %q after other player exceeded pres", president)
	}
}

// TestStockBuyWithSaleProfit checks to make sure that stocks can be bought with the money earned
// from sales made in the same turn as the purchase.
func TestStockBuyWithSaleProfit(t *testing.T) {
	playerName := "player"
	game := NewGame([]string{playerName})
	game.Players[playerName].Cash = 400

	// Chose the companies whose stocks the player will buy and sell.
	company1 := randomCompany(false)
	company2 := randomCompany(false)
	company3 := randomCompany(false)
	for company2 == company1 {
		company2 = randomCompany(false)
	}
	for company3 == company1 || company3 == company2 {
		company3 = randomCompany(false)
	}

	// First we need to buy stock in company 1 so we can sell it to buy from company 2
	if errs := startCompany(t, game, company1, 6, 60); len(errs) > 0 {
		t.Fatalf("failed to make the first purchase: %v", errs)
	}
	// Get rid of all the players remaining cash so there's no way the rest of the buys
	// will work without the money from the sale of previously bought stock.
	game.Players[playerName].Cash = 0

	turn := MarketTurn{
		Sales: []MarketAction{
			MarketAction{Company: company1, Count: 3, Price: 60},
		},
		Purchase: &MarketAction{Company: company2, Count: 3, Price: 60},
	}
	if errs := testMarketTurn(t, game, playerName, turn); len(errs) > 0 {
		t.Fatalf("failed to sell company 1 stock to buy company 2 stock: %v", errs)
	}

	turn = MarketTurn{
		Sales: []MarketAction{
			MarketAction{Company: company1, Count: 1},
			MarketAction{Company: company2, Count: 1},
		},
		Purchase: &MarketAction{Company: company3, Count: 2, Price: 60},
	}
	if errs := testMarketTurn(t, game, playerName, turn); len(errs) > 0 {
		t.Fatalf("failed to sell company 1 and company 2 stock to buy company 3 stock: %v", errs)
	}
}
