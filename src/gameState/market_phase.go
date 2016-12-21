package gameState

import (
	"fmt"

	"boardInfo"
)

func (g *Game) PerformMarketTurn(playerName string, turn MarketTurn) (errs []error) {
	defer func() {
		// If there was no error when we returned that means this action succeeded and its the
		// next player's turn.
		if len(errs) == 0 {
			g.endMarketTurn(len(turn.Sales) == 0 && turn.Purchase == nil)
		}
	}()

	if !g.marketPhase() {
		return []error{fmt.Errorf("Must be in market phase to perform market actions")}
	}

	player := g.Players[playerName]
	if player == nil {
		return []error{fmt.Errorf("No player with name %q", playerName)}
	} else if expected := g.currentTurn(); playerName != expected {
		return []error{fmt.Errorf("It is currently player %s's turn", expected)}
	}

	saleCash := 0
	for ind := range turn.Sales {
		if err := g.validateSale(player, &turn.Sales[ind]); err != nil {
			errs = append(errs, err)
		} else {
			saleCash += turn.Sales[ind].Count * turn.Sales[ind].Price
		}
	}
	if turn.Purchase != nil {
		if err := g.validateBuy(player, turn.Purchase, saleCash); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return
	}

	for _, saleInfo := range turn.Sales {
		g.sellStock(player, saleInfo)
	}
	if turn.Purchase != nil {
		g.buyStock(player, *turn.Purchase)
	}

	return nil
}

// validateAction performs the basic validation common to sales and purchases. It also fills in the
// price if the user did not provide it and the company has already been started.
func (g *Game) validateAction(action *MarketAction) (*Company, error) {
	if action.Company == "" {
		return nil, fmt.Errorf("Must specify company name to buy/sell stocks")
	}

	company := g.Companies[action.Company]
	if company == nil {
		return nil, fmt.Errorf("No company with name %q", action.Company)
	} else if action.Count == 0 {
		return nil, fmt.Errorf("Cannot buy/sell 0 stocks in %s", company.Name)
	} else if company.StockPrice != 0 {
		// This isn't really a critical error, since we could just ignore the provided price, but
		// since the user doesn't have to specify the price this is a simple way to make sure that
		// if the user claims to know the price they are in fact correct.
		if action.Price != 0 && action.Price != company.StockPrice {
			return nil, fmt.Errorf("%s shares cannot be bought/sold for $%d (price is $%d)",
				company.Name, action.Price, company.StockPrice)
		}
		action.Price = company.StockPrice
	}
	return company, nil
}

// validateBuy checks to make sure a purchase would be valid. This includes extra validation if the
// company has not yet been started. The function performs no actions on the player or the company
// as the function is primarily intended to make sure all market actions for a given turn are
// valid before any are performed.
//
// The function accepts an extra argument specifying how much the player will earn for the sale
// of other stock since this money should be available for the purchase of additional stock.
func (g *Game) validateBuy(player *Player, buyInfo *MarketAction, saleCash int) error {
	company, err := g.validateAction(buyInfo)
	if err != nil {
		return err
	}

	// If the company hasn't been started yet we need to make sure the company is allowed to be
	// started, and that the price the player wants to start it at is valid for this tech level.
	if company.StockPrice == 0 {
		if company.Restricted && g.TechLevel < 3 {
			return fmt.Errorf("%s locked until tech level 3", company.Name)
		}

		validPrice := false
		for _, option := range boardInfo.StartingStockPrices(g.TechLevel) {
			validPrice = validPrice || buyInfo.Price == option
		}
		if !validPrice {
			return fmt.Errorf("$%d not one of the valid starting prices for tech level %d",
				buyInfo.Price, g.TechLevel)
		}
	}

	// Next check to make sure the player can afford the stock they want, and that there is enough
	// non-player held stock available for them to purchase.
	if buyInfo.Count*buyInfo.Price > player.Cash+saleCash {
		return fmt.Errorf("%s has insufficient cash for %d shares of %s at $%d",
			player.Name, buyInfo.Count, company.Name, buyInfo.Price)
	} else if left := company.HeldStock + g.OrphanStocks[company.Name]; buyInfo.Count > left {
		return fmt.Errorf("%s only has %d shares remaining", company.Name, left)
	}
	return nil
}

// buyStock exchanges player cash for stock in a company. If the bank owns any stock in the company
// (orphaned stock) that must be purchased before the stock held by the company. If a player buys
// enough to acquire a controlling interest in the company the title of president is transferred to
// them.
func (g *Game) buyStock(player *Player, buyInfo MarketAction) error {
	// This function should only be called with valid MarketActions, so we shouldn't need to do
	// any error checking.
	company := g.Companies[buyInfo.Company]

	if company.StockPrice == 0 {
		company.StockPrice = buyInfo.Price
		company.PriceChange = g.timeString()
		company.BuiltTrack = []string{boardInfo.StartingLocation(company.Name)}
	}

	// We have separate variables for these instead of just using the MarketAction so we can
	// redefine how much goes to / comes from the company with the existence of orphan stocks.
	price, count := company.StockPrice*buyInfo.Count, buyInfo.Count
	player.Cash -= price
	player.Stocks[company.Name] += count

	// If this company has any orphaned stock then purchase that before we start cutting into the
	// companies holding. (No need to let the player chose, doing otherwise is throwing money away.)
	if g.OrphanStocks[company.Name] > 0 {
		if count >= g.OrphanStocks[company.Name] {
			count -= g.OrphanStocks[company.Name]
			delete(g.OrphanStocks, company.Name)
		} else {
			g.OrphanStocks[company.Name] -= count
			count = 0
		}
		price = count * company.StockPrice
	}
	company.Treasury += price
	company.HeldStock -= count

	if company.President == "" ||
		g.Players[company.President].Stocks[company.Name] < player.Stocks[company.Name] {
		company.President = player.Name
	}
	return nil
}

// validateSale checks to make sure a sale would be valid. It does not perform any actions, as this
// function is primarily intended to make sure all market actions in a given turn
func (g *Game) validateSale(player *Player, saleInfo *MarketAction) error {
	if company, err := g.validateAction(saleInfo); err != nil {
		return err
	} else if held := player.Stocks[company.Name]; held == 0 {
		return fmt.Errorf("%s has no stock in %s", player.Name, company.Name)
	} else if saleInfo.Count > held {
		return fmt.Errorf("%s only has %d shares in %s", player.Name, held, company.Name)
	} else if saleInfo.Count == held {
		// If the sum of this player's shares, the company's shares, and the orphaned shares
		// equal the total number of shares then no other player has any stock in this company
		if company.HeldStock+g.OrphanStocks[company.Name]+saleInfo.Count == 10 {
			return fmt.Errorf("Cannot sell the last player held stock in %s", company.Name)
		}
	}

	return nil
}

// sellStock exchanges player shares in a company for money from the bank. The shares then belong
// to the bank and are considered orphaned. IF the player making the sale was the president and
// they sold their majority interest the title of president is transfered.
func (g *Game) sellStock(player *Player, saleInfo MarketAction) error {
	// This function should only be called with valid MarketActions, so we shouldn't need to do
	// any error checking.
	company := g.Companies[saleInfo.Company]

	g.OrphanStocks[company.Name] += saleInfo.Count
	player.Cash += saleInfo.Count * company.StockPrice
	player.Stocks[company.Name] -= saleInfo.Count
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
	}
	return nil
}
