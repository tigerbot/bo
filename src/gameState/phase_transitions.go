package gameState

import (
	"fmt"
	"math/rand"
	"sort"

	"boardInfo"
)

func (n phaseNum) MarshalJSON() ([]byte, error) {
	result := ""
	switch n {
	case 0:
		result = `"Market"`
	case 1:
		result = `"Business 1"`
	case 2:
		result = `"Business 2"`
	}
	return []byte(result), nil
}

func (n phaseNum) Market() bool {
	return n == 0
}
func (n phaseNum) Business() bool {
	return n == 1 || n == 2
}

func (t TurnManager) Current() string {
	if len(t.Order) == 0 {
		return ""
	}
	return t.Order[t.Number%len(t.Order)]
}
func (t TurnManager) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Current())), nil
}

func (g *GlobalState) timeString() string {
	return fmt.Sprintf("%02d-%02d-%02d", g.Round, g.Phase, g.TurnManager.Number)
}

func (g *Game) beginMarketPhase() {
	g.Round += 1
	g.Phase = 0

	g.TurnManager.Number = 0
	g.TurnManager.Passes = 0
	g.TurnManager.Order = make([]string, 0, len(g.Players))
	for name, _ := range g.Players {
		g.TurnManager.Order = append(g.TurnManager.Order, name)
	}
	sort.Sort(playerSorter{list: g.TurnManager.Order, info: g.Players})
}

func (g *Game) beginBusinessPhase() {
	g.Phase += 1
	g.TurnManager.Number = 0

	g.TurnManager.Order = make([]string, 0, len(g.Companies))
	for name, company := range g.Companies {
		if company.President != "" {
			g.TurnManager.Order = append(g.TurnManager.Order, name)
		}
		company.TurnStage = ""
	}
	sort.Sort(companySorter{list: g.TurnManager.Order, info: g.Companies})
	g.Companies[g.TurnManager.Order[0]].TurnStage = "inventory"
}

func (g *Game) endMarketTurn(pass bool) {
	g.TurnManager.Number += 1

	// If every player has passed (in a row) then we are finished with the market phase and need
	// to begin the next business phase.
	if !pass {
		g.TurnManager.Passes = 0
	} else if g.TurnManager.Passes += 1; g.TurnManager.Passes == len(g.TurnManager.Order) {
		// Any companies that have orphaned stock by the end of the market phase have their stock
		// prices reduced.
		for name, _ := range g.OrphanStocks {
			company := g.Companies[name]
			company.StockPrice = boardInfo.PrevStockPrice(company.StockPrice)
		}
		g.beginBusinessPhase()
	}
}

func (g *Game) endBusinessTurn() {
	g.Companies[g.TurnManager.Current()].TurnStage = ""
	if g.TurnManager.Number += 1; g.TurnManager.Number == len(g.TurnManager.Order) {
		if g.Phase == 1 {
			g.beginBusinessPhase()
		} else {
			g.beginMarketPhase()
		}
	} else {
		g.Companies[g.TurnManager.Current()].TurnStage = "inventory"
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
	s.list[i], s.list[j] = s.list[j], s.list[i]
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
	s.list[i], s.list[j] = s.list[j], s.list[i]
}
func (s companySorter) Less(i, j int) bool {
	item1, item2 := s.info[s.list[i]], s.info[s.list[j]]
	if item1.StockPrice != item2.StockPrice {
		return item1.StockPrice > item2.StockPrice
	}
	return item1.PriceChange < item2.PriceChange
}
