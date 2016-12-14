package gameState

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

// TestMarkerTurnOrder checks to make sure the players with the least capital get to go first. It
// also make sure that the order of players that are equivalent (like at the beginning of the game)
// are sorted in random order.
func TestMarketTurnOrder(t *testing.T) {
	game := new(Game)

	game.Players = map[string]*Player{
		"1st": &Player{Cash: 250, NetWorth: 250, Name: "1st"},
		"2nd": &Player{Cash: 250, NetWorth: 250, Name: "2nd"},
		"3rd": &Player{Cash: 250, NetWorth: 250, Name: "3rd"},
		"4th": &Player{Cash: 250, NetWorth: 250, Name: "4th"},
		"5th": &Player{Cash: 250, NetWorth: 250, Name: "5th"},
		"6th": &Player{Cash: 250, NetWorth: 250, Name: "6th"},
	}
	game.beginMarketPhase()
	prevOrder := game.turnOrder
	for inc := 0; inc < 2; inc += 1 {
		game.beginMarketPhase()
		if !reflect.DeepEqual(game.turnOrder, prevOrder) {
			break
		}
	}
	if reflect.DeepEqual(game.turnOrder, prevOrder) {
		t.Errorf("player order %v produced repeatedly for even players", prevOrder)
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
		t.Errorf("company order %v doesn't match %v after presidents removed", game.turnOrder, expected)
	}
}
