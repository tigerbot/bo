package gameState

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

// TestGameStart checks various parameters of a newly created game. Probably not actually very
// useful given the current NewGame's complexity (most of it is just assignments that would
// pretty much need to be copied to be tested which just duplicates code).
func TestGameStart(t *testing.T) {
	playerCnt := 3 + rand.Intn(4)
	playerNames := make([]string, playerCnt)
	for ind := range playerNames {
		playerNames[ind] = fmt.Sprintf("%d", ind)
	}
	game := NewGame(playerNames)

	if game.Round != 1 {
		t.Errorf("new game didn't start on round 1: %v", game.Round)
	}
	if !game.marketPhase() {
		t.Errorf("new game didn't start in the market phase: %v", game.Phase)
	}

	if len(game.Players) != playerCnt {
		t.Errorf("new game has %d players, expected %d", len(game.Players), playerCnt)
	} else {
		for _, name := range playerNames {
			if game.Players[name] == nil {
				t.Errorf("new game missing player %q", name)
			}
		}
	}
	commonCash := game.Players[playerNames[0]].Cash
	for name, player := range game.Players {
		if name != player.Name {
			t.Errorf("player names %q and %q mismatch", name, player.Name)
		}
		if len(player.Stocks) != 0 {
			t.Errorf("player %q started the game with stock: %v", name, player.Stocks)
		}
		if player.Cash != commonCash {
			t.Errorf("player %q started with %d$, others have %d$", name, player.Cash, commonCash)
		}
		if player.Cash != player.NetWorth {
			t.Errorf("player %q has value mismatch: %v$ != %v$", name, player.Cash, player.NetWorth)
		}
	}

	if len(game.Companies) != 10 {
		t.Errorf("new game has %d companies, expected 10", len(game.Companies))
	}
	for name, company := range game.Companies {
		if name != company.Name {
			t.Errorf("company names %q and %q mismatch", name, company.Name)
		}
		if len(company.BuiltTrack) != 0 {
			t.Errorf("company %q started the game with tracks: %v", name, company.BuiltTrack)
		}
		if company.StockPrice != 0 {
			t.Errorf("company %q started with non-zero price %d$", name, company.StockPrice)
		}
		if company.HeldStock != 10 {
			t.Errorf("company %q started with %d shares, expected 10", name, company.HeldStock)
		}
	}
}

// TestMarkerTurnOrder checks to make sure the players with the least capital get to go first. It
// also make sure that the order of players that are equivalent (like at the beginning of the game)
// are sorted in random order.
func TestMarketTurnOrder(t *testing.T) {
	game := NewGame([]string{"1st", "2nd", "3rd", "4th", "5th", "6th"})
	prevOrder := game.turnOrder
	for inc := 0; inc < 2; inc += 1 {
		game.beginMarketPhase()
		if !reflect.DeepEqual(game.turnOrder, prevOrder) {
			break
		}
	}
	if reflect.DeepEqual(game.turnOrder, prevOrder) {
		t.Errorf("player order %v produced repeatedly for starting players", prevOrder)
	}

	game.Players = map[string]*Player{
		"1st": &Player{Cash: 100, NetWorth: 999, Name: "1st"},
		"2nd": &Player{Cash: 200, NetWorth: 150, Name: "2nd"},
		"3rd": &Player{Cash: 200, NetWorth: 200, Name: "3rd"},
		"4th": &Player{Cash: 300, NetWorth: 100, Name: "4th"},
		"5th": &Player{Cash: 400, NetWorth: 100, Name: "5th"},
		"6th": &Player{Cash: 400, NetWorth: 600, Name: "6th"},
	}
	game.beginMarketPhase()

	expected := []string{"1st", "2nd", "3rd", "4th", "5th", "6th"}
	if !reflect.DeepEqual(game.turnOrder, expected) {
		t.Errorf("initial player order %v doesn't match %v", game.turnOrder, expected)
	}
}

// TestBusinessTurnOrder checks to make sure the most valuable companies get to go first in the
// business phases. If the event of a tie of stock price the tie should be settled by which
// company had that value first.
func TestBusinessTurnOrder(t *testing.T) {
	game := new(Game)
	game.Companies = map[string]*Company{
		"1st": &Company{StockPrice: 500, PriceChange: "01-02-01", President: "b", Name: "1st"},
		"2nd": &Company{StockPrice: 500, PriceChange: "01-02-02", President: "b", Name: "2nd"},
		"3rd": &Company{StockPrice: 400, PriceChange: "01-02-00", President: "b", Name: "3rd"},
		"4th": &Company{StockPrice: 300, PriceChange: "02-00-04", President: "b", Name: "4th"},
		"5th": &Company{StockPrice: 300, PriceChange: "02-00-05", President: "b", Name: "5th"},
		"6th": &Company{StockPrice: 300, PriceChange: "02-01-02", President: "b", Name: "6th"},
		"7th": &Company{StockPrice: 200, PriceChange: "01-00-08", President: "b", Name: "7th"},
		"8th": &Company{StockPrice: 200, PriceChange: "01-01-00", President: "b", Name: "8th"},
		"9th": &Company{StockPrice: 100, PriceChange: "01-02-00", President: "b", Name: "9th"},
	}
	game.beginBusinessPhase()

	expected := []string{"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th", "9th"}
	if !reflect.DeepEqual(game.turnOrder, expected) {
		t.Errorf("initial company order %v doesn't match %v", game.turnOrder, expected)
	}

	game.Companies["7th"].President = ""
	game.Companies["8th"].President = ""
	game.beginBusinessPhase()
	expected = []string{"1st", "2nd", "3rd", "4th", "5th", "6th", "9th"}
	if !reflect.DeepEqual(game.turnOrder, expected) {
		t.Errorf("company order %v doesn't match %v after presidents removed",
			game.turnOrder, expected)
	}
}

// TestMarketPhaseEnd checks to make sure the market phase ends when all players have passed back
// to back. It also checks to make sure any companies with orphaned stocks have the stock price
// decrease at the end of the market phase.
func TestMarketPhaseEnd(t *testing.T) {
	playerCnt := 3 + rand.Intn(4)
	playerNames := make([]string, playerCnt)
	for ind := range playerNames {
		playerNames[ind] = fmt.Sprintf("%d", ind)
	}
	game := NewGame(playerNames)

	if !game.marketPhase() {
		t.Fatalf("new game didn't start in the market phase: %v", game.Phase)
	}
	if game.passes != 0 {
		t.Fatalf("new game started off with market turn passes: %v", game.passes)
	}
	if game.turn != 0 {
		t.Fatalf("new game started off with market turn passes: %v", game.passes)
	}

	// Create an orphan company to make sure its stock price is reduced only at the end.
	const origPrice = 200
	var orphanCompany *Company
	for _, company := range game.Companies {
		orphanCompany = company
		break
	}
	orphanCompany.StockPrice = origPrice
	game.OrphanStocks[orphanCompany.Name] = 2

	turn := game.turn
	endTurn := func(pass bool) {
		game.endMarketTurn(pass)
		turn += 1
		if !game.marketPhase() {
			t.Fatalf("market phase ended prematurely after turn %d with %d players",
				turn, playerCnt)
		}
		if game.turn != turn {
			t.Errorf("game turn counter %d and test turn counter %d disagree",
				game.turn, turn)
			turn = game.turn
		}
		if orphanCompany.StockPrice != origPrice {
			t.Errorf("orphan company stock changed prematurely: %d -> %d",
				origPrice, orphanCompany.StockPrice)
		}
	}
	for loopInd := 0; loopInd < 10; loopInd += 1 {
		for passes := rand.Intn(playerCnt); passes > 0; passes -= 1 {
			endTurn(true)
		}
		for actions := rand.Intn(playerCnt) + 1; actions > 0; actions -= 1 {
			endTurn(false)
		}
	}

	// Now actually fulfill the market phase end requirements.
	for passes := 0; passes < playerCnt-1; passes += 1 {
		endTurn(true)
	}
	game.endMarketTurn(true)
	if game.marketPhase() {
		t.Fatal("market phase did not end after every player passed")
	}
	if orphanCompany.StockPrice >= origPrice {
		t.Errorf("orphan stock price %d did not drop from %d", orphanCompany.StockPrice, origPrice)
	}
}

// TestBusinessPhaseEnd checks to make sure each business phase ends, when all companies with a
// president have taken their turn. It also checks to make sure the first business phase ending
// means the beginning of the second business phase, and that the second business phase ending
// means the beginning of the market phase of the next round.
func TestBusinessPhaseEnd(t *testing.T) {
	// Start off with a game in the market phase so we can call beginBusinessPhase.
	game := &Game{
		GlobalState: GlobalState{Round: 1, Phase: 0},
		Companies: map[string]*Company{
			"1st": {Name: "1st", President: "president", StockPrice: 350},
			"2nd": {Name: "2nd", President: "president", StockPrice: 300},
			"3rd": {Name: "3rd", President: "president", StockPrice: 250},
			"4th": {Name: "4th", President: "", StockPrice: 200},
			"5th": {Name: "5th", President: "president", StockPrice: 0},
			"6th": {Name: "6th", President: "", StockPrice: 0},
		},
	}
	if !game.marketPhase() {
		t.Fatal("expected game to start in market phase, but it didn't")
	}

	game.beginBusinessPhase()
	if !game.businessPhase() {
		t.Fatalf("did not enter business phase after calling beginBusinessPhase: %v", game.Phase)
	} else if len(game.turnOrder) != 4 {
		t.Fatalf("expected 4 items in the turn order: %v", game.turnOrder)
	}

	checkTurnEndState := func(callNum, round, phase int, business bool) {
		game.endBusinessTurn()
		if business != game.businessPhase() {
			t.Errorf("expected businessPhase to be %v after %d calls", business, callNum)
		} else if phase != game.Phase {
			t.Errorf("expected phase %d after %d calls, found %d", phase, callNum, game.Phase)
		} else if round != game.Round {
			t.Errorf("expected round %d after %d calls, found %d", round, callNum, game.Round)
		}
	}

	checkTurnEndState(1, 1, 1, true)
	checkTurnEndState(2, 1, 1, true)
	checkTurnEndState(3, 1, 1, true)
	checkTurnEndState(4, 1, 2, true)
	checkTurnEndState(5, 1, 2, true)
	checkTurnEndState(6, 1, 2, true)
	checkTurnEndState(7, 1, 2, true)
	checkTurnEndState(8, 2, 0, false)
}
