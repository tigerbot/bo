package gameState

import (
	"math/rand"
	"testing"
)

// startingPrices is the list of valid stock prices for each tech level (as printed on the physical
// board of the original game).
var startingPrices = [][3]int{
	[3]int{55, 60, 66},
	[3]int{60, 66, 74},
	[3]int{66, 74, 82},
	[3]int{74, 82, 91},
	[3]int{82, 91, 100},
}

// startCompany is a convenience function for the tests that check starting companies under various
// conditions. It creates a new game for every call, initializes some of the specified parameters,
// then performs the market action with the remaining parameters.
func startCompany(t *testing.T, game *Game, company string, count, price int) []error {
	action := MarketAction{
		Company: company,
		Count:   count,
		Price:   price,
	}

	playerName := game.currentTurn()
	cash := game.Players[playerName].Cash

	if errs := testMarketTurn(t, game, playerName, MarketTurn{Purchase: &action}); len(errs) > 0 {
		return errs
	}

	// Test the things that should be the same no matter what the initial conditions are.
	if endPrice := game.Companies[company].StockPrice; endPrice != price {
		t.Errorf("after starting %s at $%d, stock price is $%d", company, price, endPrice)
	}
	if stockLeft := game.Companies[company].HeldStock; stockLeft != 10-count {
		t.Errorf("%d stock left after starting purchase of %d", stockLeft, count)
	}
	if president := game.Companies[company].President; president != playerName {
		t.Errorf("%s started company, but company president is %q", playerName, president)
	}
	if playerCash := game.Players[playerName].Cash; playerCash != cash-count*price {
		t.Errorf("player started with $%d, bought %d shares at $%d, left with $%d",
			cash, count, price, playerCash)
	}
	if playerStock := game.Players[playerName].Stocks[company]; playerStock != count {
		t.Errorf("player has %d shares of %d after buying %d", playerStock, count)
	}

	return nil
}

func startCompanyNewGame(t *testing.T, company string, count, price, cash, techLvl int) []error {
	game := NewGame([]string{"1st", "2nd", "3rd"})
	game.TechLevel = techLvl
	playerName := game.TurnOrder[0]
	game.Players[playerName].Cash = cash

	return startCompany(t, game, company, count, price)
}

// TestCompanyStart checks to make sure a company cannot be started by a player with insufficient
// cash, or by purchasing more than 100% of the company's stock. It also makes sure that companies
// can only be started at valid initial stock prices.
func TestCompanyStart(t *testing.T) {
	company := randomCompany(false)

	for ind, prices := range startingPrices {
		techLvl := ind + 1

		for _, price := range prices {
			if errs := startCompanyNewGame(t, company, 4, price, 500, techLvl); len(errs) > 0 {
				t.Errorf("starting %s in tech level %d at $%d failed: %v",
					company, techLvl, price, errs)
			}
		}

		invalidPrices := []int{
			prices[0] - (rand.Intn(20) + 1),
			prices[0] + rand.Intn(prices[1]-prices[0]-1) + 1,
			prices[1] + rand.Intn(prices[2]-prices[1]-1) + 1,
			prices[2] + rand.Intn(20) + 1,
		}
		for _, price := range invalidPrices {
			if errs := startCompanyNewGame(t, company, 1, price, price*2, techLvl); len(errs) == 0 {
				t.Errorf("starting %s in tech level %d at $%d (invalid) did not error",
					company, techLvl, price)
			}
		}
	}

	price := startingPrices[0][1]

	count := rand.Intn(9) + 1
	if errs := startCompanyNewGame(t, company, count, price, count*price-20, 1); len(errs) == 0 {
		t.Errorf("attempt to buy %d %s shares at $%d with insufficient cash did not error",
			count, company, price)
	}

	count = 11
	if errs := startCompanyNewGame(t, company, count, price, count*price, 1); len(errs) == 0 {
		t.Errorf("attempt to buy %d %s shares at $%d with insufficient cash did not error",
			count, company, price)
	}
}

// TestCompanyStartTech3 makes sure the companies labeled at only available after tech level 3
// cannot be started before that.
func TestCompanyStartTech3(t *testing.T) {
	company := randomCompany(true)

	for ind, prices := range startingPrices[:2] {
		techLvl := ind + 1
		price := prices[rand.Intn(3)]
		if errs := startCompanyNewGame(t, company, 4, price, 500, techLvl); len(errs) == 0 {
			t.Errorf("starting %s during tech level %d did not error", company, techLvl)
		}
	}

	for ind, prices := range startingPrices[2:] {
		techLvl := ind + 3
		price := prices[rand.Intn(3)]
		if errs := startCompanyNewGame(t, company, 4, price, 500, techLvl); len(errs) > 0 {
			t.Errorf("starting %s during tech level %d failed: %v", company, techLvl, errs)
		}
	}
}
