package gameState

import (
	"fmt"
	"sort"

	"boardInfo"
)

func (g *Game) HandleCompanyEarnings(playerName string, earnings CompanyEarnings) []error {
	if !g.Phase.Business() {
		return []error{fmt.Errorf("Must be in a business phase to perform business actions")}
	}
	company := g.Companies[g.TurnManager.Current()]
	if playerName != company.President {
		return []error{
			fmt.Errorf("It's %s's turn and %s is the president", company.Name, company.President)}
	}
	if g.Stage != "earnings" {
		return []error{fmt.Errorf("%s is not ready to handle its earnings", company.Name)}
	}

	if len(earnings.Serviced) == 0 {
		cities := boardInfo.Cities(company.BuiltTrack...)
		capacity := 0
		for ind, count := range company.Equipment {
			capacity += (ind + 1) * count
		}

		if len(cities) > capacity {
			boardInfo.SortCities(cities, g.TechLevel)
			cities = cities[:capacity]
		}

		earnings.Serviced = make([]string, len(cities))
		for ind, city := range cities {
			earnings.Serviced[ind] = city.Location
		}
	}

	if errs := g.validateServicedCities(company, earnings); len(errs) > 0 {
		return errs
	}

	costs := 0
	for _, count := range company.Equipment {
		costs += 10 * g.TechLevel * count
	}

	gross := 0
	// The company receives income from mined coal if it has any capital equipment.
	if costs > 0 {
		gross += 40 * company.CoalMined
	}
	// And the company receives income from every city serviced dependent on tech level.
	for _, city := range boardInfo.Cities(earnings.Serviced...) {
		gross += city.Revenue[g.TechLevel-1]
	}

	net := gross - costs
	defer func(prevPrice int) {
		company.NetIncome = net

		// If the stock price for this company has changed then update the net worth for
		// any players that have stock in this company.
		if company.StockPrice != prevPrice {
			for _, player := range g.Players {
				if player.Stocks[company.Name] > 0 {
					player.NetWorth = player.Cash
					for name, count := range player.Stocks {
						player.NetWorth += count * g.Companies[name].StockPrice
					}
				}
			}
		}
	}(company.StockPrice)

	if net <= 0 {
		// A lot happens for unprofitable companies, so it was put into a different function
		g.handleNonprofitable(company)
		net = 0
	} else if net < company.NetIncome {
		company.StockPrice = boardInfo.PrevStockPrice(company.StockPrice)
		company.PriceChange = g.timeString()
	} else if net > company.NetIncome && earnings.Dividends {
		company.StockPrice = boardInfo.NextStockPrice(company.StockPrice)
		company.PriceChange = g.timeString()
	}

	if !earnings.Dividends {
		company.Treasury += net
	} else {
		perShare := net / 10
		company.Treasury += perShare * company.HeldStock
		for _, player := range g.Players {
			total := perShare * player.Stocks[company.Name]
			player.Cash += total
			player.NetWorth += total
		}
	}

	g.endBusinessTurn()
	return nil
}

// validateServicedCities makes sure the company has the capability of servicing all the cities
// the president indicated should be covered. It realistically doesn't need to be attached to
// the game object, but since it is basically alone in this regard we keep things consistent.
func (g *Game) validateServicedCities(company *Company, earnings CompanyEarnings) []error {
	var errs []error

	// First we need to make sure the company isn't trying to service more cities than it has
	// the ability to.
	capacity := 0
	for ind, count := range company.Equipment {
		capacity += (ind + 1) * count
	}
	if len(earnings.Serviced) > capacity {
		errs = append(errs, fmt.Errorf("%s can only service %d cities", company.Name, capacity))
	}

	// Make sure every coordinate provided actually has a city on it.
	cities := boardInfo.Cities(earnings.Serviced...)
	if len(cities) != len(earnings.Serviced) {
		errs = append(errs, fmt.Errorf("not all map locations provided contain cities"))
	}

	// Then make sure the company isn't trying to service any cities that it doesn't have
	// any tracks in.
	for _, city := range cities {
		ind := sort.SearchStrings(company.BuiltTrack, city.Location)
		if ind < 0 || ind >= len(company.BuiltTrack) || company.BuiltTrack[ind] != city.Location {
			errs = append(errs, fmt.Errorf("%s is not present in %s to be able to service it",
				company.Name, city.Name))
		}
	}

	return errs
}

func (g *Game) handleNonprofitable(company *Company) {
	// The stock price of an unprofitable company goes back two spaces.
	company.StockPrice = boardInfo.PrevStockPrice(company.StockPrice)
	company.StockPrice = boardInfo.PrevStockPrice(company.StockPrice)
	company.PriceChange = g.timeString()

	// The president of the company looses one share to held stock for no pay
	president := g.Players[company.President]
	company.HeldStock += 1
	president.Stocks[company.Name] -= 1
	if president.Stocks[company.Name] == 0 {
		delete(president.Stocks, company.Name)
	}

	// After checking for a new president, if no one holds any stock in this company it enters
	// receivership, losing its treasury and all capital equipment, recovering its orphaned
	// stock, and setting its stock price at 50.
	for _, player := range g.Players {
		if player.Stocks[company.Name] > president.Stocks[company.Name] {
			president = player
		}
	}
	if president.Stocks[company.Name] > 0 {
		company.President = president.Name
		return
	}

	delete(g.OrphanStocks, company.Name)
	company.HeldStock = 10
	company.Treasury = 0
	for ind := range company.Equipment {
		company.Equipment[ind] = 0
	}
	company.StockPrice = 50
}
