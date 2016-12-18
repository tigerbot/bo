package gameState

import (
	"math/rand"
	"testing"
)

func sellStock(t *testing.T, game *Game, playerName, companyName string, count int) error {
	startPrice := game.Companies[companyName].StockPrice
	startHeld := game.Companies[companyName].HeldStock
	startTreasure := game.Companies[companyName].Treasury
	startOrphan := game.OrphanStocks[companyName]

	startStock := game.Players[playerName].Stocks[companyName]
	startCash := game.Players[playerName].Cash

	err := game.PerformMarketAction(playerName, MarketAction{Company: companyName, Count: -count})
	if stockPrice := game.Companies[companyName].StockPrice; stockPrice != startPrice {
		t.Errorf("selling stock changed price from $%d to $%d", startPrice, stockPrice)
	}
	if held := game.Companies[companyName].HeldStock; held != startHeld {
		t.Errorf("selling stock changed company held stock from %d to %d", startHeld, held)
	}
	if treasure := game.Companies[companyName].Treasury; treasure != startTreasure {
		t.Errorf("selling stock changed company treasury from $%d to $%d", startTreasure, treasure)
	}
	if err == nil {
		totalCost := count * startPrice
		if change := startStock - game.Players[playerName].Stocks[companyName]; change != count {
			t.Errorf("player lost %d stock after selling %d", change, count)
		}
		if change := game.Players[playerName].Cash - startCash; change != totalCost {
			t.Errorf("player made $%d after selling %d stock at $%d", change, count, startPrice)
		}
		if change := game.OrphanStocks[companyName] - startOrphan; change != count {
			t.Errorf("company orphaned stock changed by %d after player sold %d", change, count)
		}
	} else {
		if change := startStock - game.Players[playerName].Stocks[companyName]; change != 0 {
			t.Errorf("player lost %d stock after failing to sell %d", change, count)
		}
		if change := game.Players[playerName].Cash - startCash; change != 0 {
			t.Errorf("player made $%d after failing to sell %d stock at $%d",
				change, count, startPrice)
		}
		if change := game.OrphanStocks[companyName] - startOrphan; change != 0 {
			t.Errorf("company orphaned stock changed by %d after player failed to sold %d",
				change, count)
		}
	}
	return err
}

// TestStockSell checks to make sure players cannot sell more stock than they have, and that the
// money and stock changes the way it should.
func TestStockSell(t *testing.T) {
	playerName := "player"
	game := NewGame([]string{playerName})
	action := MarketAction{
		Company: randomCompany(false),
		Price:   startingPrices[0][1],
		Count:   rand.Intn(5) + 2,
	}
	game.Players[playerName].Cash = action.Price * action.Count

	if err := sellStock(t, game, playerName, action.Company, 1); err == nil {
		t.Fatal("player selling stock they do not have did not fail")
	}

	if err := game.PerformMarketAction(playerName, action); err != nil {
		t.Fatalf("player failed to acquire stock to sell: %v", err)
	}

	if err := sellStock(t, game, playerName, action.Company, action.Count+1); err == nil {
		t.Fatal("player selling more stock than they hold did not fail")
	}

	if err := sellStock(t, game, playerName, action.Company, action.Count); err != nil {
		t.Fatalf("player selling all their stock failed: %v", err)
	} else {
		if _, present := game.Players[playerName].Stocks[action.Company]; present {
			t.Error("selling all player stock did not remove item from stock map")
		}
		if president := game.Companies[action.Company].President; president != "" {
			t.Errorf("selling all player %s's stock left president as %q", playerName, president)
		}
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

	game.turnOrder = []string{"other"}
	if err := game.PerformMarketAction("other", action); err != nil {
		t.Fatalf("initial stock purchase failed: %v", err)
	}

	game.turnOrder = []string{"pres"}
	if err := buyStock(t, game, "pres", action.Company, 6); err != nil {
		t.Fatalf("player pres failed to purchase remaining stock: %v", err)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("player pres buying remaining stock left president as %q", president)
	}

	if err := sellStock(t, game, "pres", action.Company, 2); err != nil {
		t.Fatalf("pres failed to sell part of their stock: %v", err)
	} else if president := game.Companies[action.Company].President; president != "pres" {
		t.Fatalf("selling enough stock to tie other player changed president to %q", president)
	}

	if err := sellStock(t, game, "pres", action.Company, 2); err != nil {
		t.Fatalf("pres failed to sell part of their stock: %v", err)
	} else if president := game.Companies[action.Company].President; president != "other" {
		t.Fatalf("selling enough stock to lose to other player left president as %q", president)
	}
}
