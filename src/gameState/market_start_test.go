package gameState

import (
	"math/rand"
	"reflect"
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

	playerName := game.TurnManager.Current()
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
	game.Players[game.TurnManager.Current()].Cash = cash

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
		t.Errorf("attempt to buy %d %s shares did not error",
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

// TestStartedTurnOrder makes sure that company turn order in the first business round is the
// same as the order they were purchased in if they at start out with the same price. The partly
// duplicates TestBusinessTurnOrder, but the main purpose of this test is really to make sure the
// parameters needed to properly sort the list are set when companies are started.
func TestStartedTurnOrder(t *testing.T) {
	game := NewGame([]string{"1st", "2nd", "3rd", "4th"})

	companyList := make([]string, 0, len(companyInitCond))
	for name, start := range companyInitCond {
		if !start.tech3 {
			companyList = append(companyList, name)
		}
	}
	for ind := range companyList {
		swap := rand.Intn(ind + 1)
		companyList[ind], companyList[swap] = companyList[swap], companyList[ind]
	}
	companyList = companyList[:4]

	// Now that we've determined which companies to start and in which order we need to actually
	// start them, then pass to end the market phase.
	for ind, companyName := range companyList {
		if errs := startCompany(t, game, companyName, 5, startingPrices[0][1]); len(errs) > 0 {
			t.Fatalf("failed to start company #%d (%s): %v", ind, companyName, errs)
		}
	}
	for _ = range companyList {
		if errs := game.PerformMarketTurn(game.TurnManager.Current(), MarketTurn{}); len(errs) > 0 {
			t.Fatalf("failed to pass market turn to advance phase: %v", errs)
		}
	}
	if !game.Phase.Business() {
		t.Fatal("failed to enter business phase")
	}
	if !reflect.DeepEqual(companyList, game.TurnManager.Order) {
		t.Errorf("company turn order %v != starting order %v", game.TurnManager.Order, companyList)
	}
}
