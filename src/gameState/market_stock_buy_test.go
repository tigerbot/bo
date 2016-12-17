package gameState

import (
	"math/rand"
	"testing"
)

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
