package gameState

import (
	"fmt"
	"math/rand"
	"sort"
)

func (g *GlobalState) timeString() string {
	result := fmt.Sprintf("%02d-%02d-%02d", g.Round, g.Phase, g.turn)
	return result
}

func (g *GlobalState) marketPhase() bool {
	return g.Phase == 0
}
func (g *GlobalState) businessPhase() bool {
	return g.Phase == 1 || g.Phase == 2
}

func (g *GlobalState) currentTurn() string {
	return g.turnOrder[g.turn%len(g.turnOrder)]
}

func (g *Game) beginMarketPhase() {
	g.Round += 1
	g.Phase = 0
	g.turn = 0
	g.passes = 0

	g.turnOrder = make([]string, 0, len(g.Players))
	for name, _ := range g.Players {
		g.turnOrder = append(g.turnOrder, name)
	}
	sort.Sort(playerSorter{list: g.turnOrder, info: g.Players})
}

func (g *Game) beginBusinessPhase() {
	g.Phase += 1
	g.turn = 0

	g.turnOrder = make([]string, 0, len(g.Companies))
	for name, company := range g.Companies {
		if company.President != "" {
			g.turnOrder = append(g.turnOrder, name)
		}
	}
	sort.Sort(companySorter{list: g.turnOrder, info: g.Companies})
}

func (g *Game) endMarketTurn(pass bool) {
	g.turn += 1

	// If every player has passed (in a row) then we are finished with the market phase and need
	// to begin the next business phase.
	if !pass {
		g.passes = 0
	} else if g.passes += 1; g.passes == len(g.turnOrder) {
		g.beginBusinessPhase()
	}
}

func (g *Game) endBusinessTurn() {
	if g.turn += 1; g.turn == len(g.turnOrder) {
		if g.Phase == 1 {
			g.beginBusinessPhase()
		} else {
			g.beginMarketPhase()
		}
	}
}

// The sorters implement sort.Interface and allow us to sort lists of the player and company names
// to determine turn order for the next phase.
type (
	playerSorter struct {
		list []string
		info map[string]*Player
	}
	companySorter struct {
		list []string
		info map[string]*Company
	}
)

func (s playerSorter) Len() int {
	return len(s.list)
}
func (s playerSorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[i], s.list[j]
}
func (s playerSorter) Less(i, j int) bool {
	item1, item2 := s.info[s.list[i]], s.info[s.list[j]]
	if item1.Cash != item2.Cash {
		return item1.Cash < item2.Cash
	} else if item1.NetWorth != item2.NetWorth {
		return item1.NetWorth < item2.NetWorth
	}
	return rand.Float32() < 0.5
}

func (s companySorter) Len() int {
	return len(s.list)
}
func (s companySorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[i], s.list[j]
}
func (s companySorter) Less(i, j int) bool {
	item1, item2 := s.info[s.list[i]], s.info[s.list[j]]
	if item1.StockPrice != item2.StockPrice {
		return item1.StockPrice > item2.StockPrice
	}
	return item1.PriceChange < item2.PriceChange
}
