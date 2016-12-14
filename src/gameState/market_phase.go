package gameState

import (
	"fmt"

	"boardInfo"
)

type MarketAction struct {
	Company string `json:"company,omitempty"`
	Count   int    `json:"count,omitempty"`
	Price   int    `json:"price,omitempty"`
}

func (g *Game) PerformMarketAction(playerName string, action MarketAction) (err error) {
	defer func() {
		// If there was no error when we returned that means this action succeeded and its the
		// next player's turn.
		if err != nil {
			g.endMarketTurn(action.Count == 0)
		}
	}()

	if !g.marketPhase() {
		return fmt.Errorf("Must be in market phase to perform market actions")
	}

	player := g.Players[playerName]
	if player == nil {
		return fmt.Errorf("No player with name %q", playerName)
	} else if expected := g.currentTurn(); playerName != expected {
		return fmt.Errorf("It is currently player %s's turn", expected)
	}

	if action.Company == "" || action.Count == 0 {
		if action.Count != 0 {
			return fmt.Errorf("Must specify company name to buy/sell stocks")
		} else if action.Company != "" {
			return fmt.Errorf("Cannot buy/sell 0 stocks in %s", action.Company)
		}
		return nil
	}

	company := g.Companies[action.Company]
	if company == nil {
		return fmt.Errorf("No company with name %q", action.Company)
	}

	// Just an error check if the user provided a price for a started company, make sure it's
	// the right price.
	if action.Price != 0 && company.StockPrice != 0 && action.Price != company.StockPrice {
		return fmt.Errorf("%s cannot be bought/sold for %d per share", action.Company, action.Price)
	}

	// Use the sign of the count and the current state of the company to determine the action
	if action.Count < 0 {
		return g.sellStock(player, company, -action.Count)
	}
	if company.StockPrice == 0 {
		return g.startCompany(player, company, action.Count, action.Price)
	}
	return g.buyStock(player, company, action.Count)
}

func (g *Game) startCompany(player *Player, company *Company, count, price int) error {
	if company.Restricted && g.TechLevel < 3 {
		return fmt.Errorf("%s locked until tech level 3", company.Name)
	}

	validPrice := false
	for _, option := range boardInfo.StartingStockPrices(g.TechLevel) {
		validPrice = validPrice || price == option
	}
	if !validPrice {
		return fmt.Errorf("%d invalid starting price for tech level %d", price, g.TechLevel)
	}

	// Initial the price, then re-use the code in the buy function, clearing the stock price
	// if that transaction failed for any reason to allow the price to be set later.
	company.StockPrice = price
	if err := g.buyStock(player, company, count); err != nil {
		company.StockPrice = 0
		return err
	}

	company.PriceChange = g.timeString()
	return nil
}

func (g *Game) buyStock(player *Player, company *Company, count int) error {
	price := count * company.StockPrice
	if count > company.HeldStock {
		return fmt.Errorf("%s only has %d shares remaining", company.Name, company.HeldStock)
	} else if price > player.Cash {
		return fmt.Errorf("%s has insufficient cash for %d shares of %s",
			player.Name, count, company.Name)
	}
	player.Cash -= price
	company.Treasury += price

	company.HeldStock -= count
	player.Stocks[company.Name] += count
	if company.President == "" ||
		g.Players[company.President].Stocks[company.Name] < player.Stocks[company.Name] {
		company.President = player.Name
	}
	return nil
}

func (g *Game) sellStock(player *Player, company *Company, count int) error {
	if held := player.Stocks[company.Name]; held == 0 {
		return fmt.Errorf("%s has no stock in %s", player.Name, company.Name)
	} else if count > held {
		return fmt.Errorf("%s only has %d shares in %s", player.Name, held, company.Name)
	}
	player.Cash += count * company.StockPrice

	g.OrphanStocks[company.Name] += count
	player.Stocks[company.Name] -= count
	if player.Stocks[company.Name] == 0 {
		delete(player.Stocks, company.Name)
	}

	if company.President == player.Name {
		president := player
		for otherName, otherPlayer := range g.Players {
			if otherPlayer.Stocks[company.Name] > president.Stocks[company.Name] {
				company.President = otherName
				president = otherPlayer
			}
		}
		if president.Stocks[company.Name] == 0 {
			company.President = ""
		}
	}
	return nil
}
