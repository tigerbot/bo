package gameState

import (
	"math/rand"
	"testing"
)

func sellStock(t *testing.T, game *Game, playerName, companyName string, count int) []error {
	startPrice := game.Companies[companyName].StockPrice
	startHeld := game.Companies[companyName].HeldStock
	startTreasure := game.Companies[companyName].Treasury
	startOrphan := game.OrphanStocks[companyName]

	startStock := game.Players[playerName].Stocks[companyName]
	startCash := game.Players[playerName].Cash
	totalCost := count * startPrice

	turn := MarketTurn{
		Sales: []MarketAction{
			MarketAction{Company: companyName, Count: count},
		},
	}
	if errs := testMarketTurn(t, game, playerName, turn); len(errs) > 0 {
		return errs
	}

	if stockPrice := game.Companies[companyName].StockPrice; stockPrice != startPrice {
		t.Errorf("selling stock changed price from $%d to $%d", startPrice, stockPrice)
	}
	if held := game.Companies[companyName].HeldStock; held != startHeld {
		t.Errorf("selling stock changed company held stock from %d to %d", startHeld, held)
	}
	if treasure := game.Companies[companyName].Treasury; treasure != startTreasure {
		t.Errorf("selling stock changed company treasury from $%d to $%d", startTreasure, treasure)
	}
	if change := startStock - game.Players[playerName].Stocks[companyName]; change != count {
		t.Errorf("player lost %d stock after selling %d", change, count)
	}
	if change := game.Players[playerName].Cash - startCash; change != totalCost {
		t.Errorf("player made $%d after selling %d stock at $%d", change, count, startPrice)
	}
	if change := game.OrphanStocks[companyName] - startOrphan; change != count {
		t.Errorf("company orphaned stock changed by %d after player sold %d", change, count)
	}
	return nil
}

// TestStockSell checks to make sure players cannot sell more stock than they have or the last
// player held stock.
func TestStockSell(t *testing.T) {
	playerName := "player"
	game := NewGame([]string{playerName})
	action := MarketAction{
		Company: randomCompany(false),
		Price:   startingPrices[0][1],
		Count:   rand.Intn(5) + 2,
	}
	game.Players[playerName].Cash = action.Price * action.Count

	if errs := sellStock(t, game, playerName, action.Company, 1); len(errs) == 0 {
		t.Fatal("player selling stock they do not have did not fail")
	}

	if err := startCompany(t, game, action.Company, action.Count, action.Price); err != nil {
		t.Fatalf("player failed to acquire stock to sell: %v", err)
	}

	if errs := sellStock(t, game, playerName, action.Company, action.Count+1); len(errs) == 0 {
		t.Fatal("player selling more stock than they hold did not fail")
	}
	if errs := sellStock(t, game, playerName, action.Company, action.Count); len(errs) == 0 {
		t.Fatal("player selling all their stock when none held by other players did not fail")
	}
	if errs := sellStock(t, game, playerName, action.Company, action.Count-1); len(errs) > 0 {
		t.Fatalf("player failed selling all but their last stock: %v", errs)
	}

	// Fabricate a player to have a share in the company so we can sell the last one.
	game.Players["fake"] = &Player{Name: "fake", Stocks: map[string]int{action.Company: 1}}
	game.OrphanStocks[action.Company] -= 1
	if errs := sellStock(t, game, playerName, action.Company, 1); len(errs) > 0 {
		t.Fatalf("player failed selling their last stock after other play acquire some: %v", errs)
	} else if _, present := game.Players[playerName].Stocks[action.Company]; present {
		t.Error("selling last of player's shares didn't remove entry from map")
	}
}

// TestStockSellStepdown checks to make sure that if a company president sells enough stock that
// another player holds a majority of the company the title is transferred.
func TestStockSellStepdown(t *testing.T) {
	game := NewGame([]string{"pres", "other"})
	action := MarketAction{
		Company: randomCompany(false),
		Price:   startingPrices[0][1],
		Count:   4,
	}
	game.Players["pres"].Cash = 6 * action.Price
	game.Players["other"].Cash = 4 * action.Price

	game.TurnManager.Order = []string{"other"}
	if errs := startCompany(t, game, action.Company, action.Count, action.Price); len(errs) > 0 {
		t.Fatalf("initial stock purchase failed: %v", errs)
	}

	game.TurnManager.Order = []string{"pres"}
	if errs := buyStock(t, game, "pres", action.Company, 6); len(errs) > 0 {
		t.Fatalf("player pres failed to purchase remaining stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("player pres buying remaining stock left president as %q", president)
	}

	if errs := sellStock(t, game, "pres", action.Company, 2); len(errs) > 0 {
		t.Fatalf("pres failed to sell part of their stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("selling enough stock to tie other player changed president to %q", president)
	}

	if errs := sellStock(t, game, "pres", action.Company, 2); len(errs) > 0 {
		t.Fatalf("pres failed to sell part of their stock: %v", errs)
	} else if president := game.Companies[action.Company].President; president != "other" {
		t.Fatalf("selling enough stock to lose to other player left president as %q", president)
	}
}
