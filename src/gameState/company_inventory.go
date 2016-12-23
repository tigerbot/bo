package gameState

import (
	"fmt"
	"sort"

	"boardInfo"
)

func stringInSlice(value string, slice []string) bool {
	for _, content := range slice {
		if value == content {
			return true
		}
	}
	return false
}

func (g *Game) UpdateCompanyInventory(playerName string, update CompanyInventory) []error {
	if !g.Phase.Business() {
		return []error{fmt.Errorf("Must be in a business phase to perform business actions")}
	}
	company := g.Companies[g.TurnManager.Current()]

	if playerName != company.President {
		return []error{
			fmt.Errorf("It's %s's turn and %s is the president", company.Name, company.President),
		}
	}
	if company.TurnStage != "inventory" {
		return []error{fmt.Errorf("%s has already updated its inventory", company.Name)}
	}

	var errs []error
	errs = append(errs, g.validateBusinessExpense(company, update)...)
	errs = append(errs, g.validateBuildLimits(company, update)...)
	errs = append(errs, g.validateCityRestrictions(company, update)...)
	if len(errs) > 0 {
		return errs
	}

	for ind, count := range update.Scrap {
		techLvl := ind + 1
		company.Treasury += 20 * techLvl * count
		company.Equipment[ind] -= count
	}
	for ind := 0; ind < update.Buy; ind += 1 {
		g.TrainsBought += 1
		g.TechLevel = boardInfo.TechLevel(g.TrainsBought)
		company.Treasury -= boardInfo.TrainCost(g.TrainsBought)
		company.Equipment[g.TechLevel] += 1
	}
	if update.Coal != "" {
		company.CoalMined += 1
		for ind, coord := range g.UnminedCoal {
			if coord == update.Coal {
				g.UnminedCoal = append(g.UnminedCoal[:ind], g.UnminedCoal[ind+1:]...)
			}
		}
	}
	company.BuiltTrack = append(company.BuiltTrack, update.Track...)
	company.UnbuiltTrack -= len(update.Track)
	sort.Strings(company.BuiltTrack)

	company.TurnStage = "earnings"
	return nil
}

func (g *Game) validateBusinessExpense(company *Company, update CompanyInventory) []error {
	var errs []error

	availableMoney := company.Treasury
	for ind, count := range update.Scrap {
		techLvl := ind + 1
		if count > company.Equipment[ind] {
			count = company.Equipment[ind]
			errs = append(errs, fmt.Errorf("%s only has %d tech level %d equipment to scrap",
				company.Name, company.Equipment[ind], techLvl))
		}
		availableMoney += 20 * techLvl * count
	}

	equipCost := 0
	for ind := 1; ind <= update.Buy; ind += 1 {
		equipCost += boardInfo.TrainCost(g.TrainsBought + ind)
	}
	trackCost := 0
	for _, hexCoord := range update.Track {
		if cost := boardInfo.BuildCost(hexCoord); cost > 0 {
			trackCost += cost
		} else {
			errs = append(errs, fmt.Errorf("Invalid hex coordinate %q to build track", hexCoord))
		}
	}

	if equipCost+trackCost > availableMoney {
		errs = append(errs, fmt.Errorf("insufficient money to perform all action"))
	}
	return errs
}

func (g *Game) validateBuildLimits(company *Company, update CompanyInventory) []error {
	var errs []error

	if len(update.Track) > 0 && update.Coal != "" {
		errs = append(errs, fmt.Errorf("cannot build track and mine coal on the same turn"))
	}
	if update.Coal != "" {
		if !stringInSlice(update.Coal, g.UnminedCoal) {
			errs = append(errs, fmt.Errorf("no coal located at %q to be mined", update.Coal))
		}
		if !stringInSlice(update.Coal, company.BuiltTrack) {
			errs = append(errs, fmt.Errorf("%s has no track on %q to allow mining",
				company.Name, update.Coal))
		}
	}
	if len(update.Track) > company.UnbuiltTrack {
		errs = append(errs, fmt.Errorf("%s only has %d unbuilt tracks remaining",
			company.Name, company.UnbuiltTrack))
	}
	if techLvl := boardInfo.TechLevel(g.TrainsBought + update.Buy); len(update.Track) > techLvl {
		errs = append(errs, fmt.Errorf("cannot build more than %d track per turn", techLvl))
	}
	for _, coord := range update.Track {
		if stringInSlice(coord, company.BuiltTrack) {
			errs = append(errs, fmt.Errorf("%s already has track built on %q", company.Name, coord))
		}
	}
	if !boardInfo.TilesContiguous(company.BuiltTrack, update.Track) {
		errs = append(errs, fmt.Errorf("all built track must be connected"))
	}

	return errs
}

func (g *Game) validateCityRestrictions(company *Company, update CompanyInventory) []error {
	var errs []error

	// The number of railroads allow in each city is the same as the tech level expect on tech
	// level one, where the limit is 2
	cityCapacity := boardInfo.TechLevel(g.TrainsBought + update.Buy)
	if cityCapacity < 2 {
		cityCapacity = 2
	}

	for _, city := range boardInfo.Cities(update.Track...) {
		// TODO: make sure no company builds in Pennsylvania's cities before it does.

		// Check to make sure there is still enough space in the city for another railroad.
		if city.Exception != "universal" {
			existent := 0
			for _, c := range g.Companies {
				if stringInSlice(city.Location, c.BuiltTrack) {
					existent += 1
				}
			}
			if existent >= cityCapacity {
				errs = append(errs, fmt.Errorf("%s already has %d railroads", city.Name, existent))
			}
		}
	}

	// TODO: make sure Pennsylvania and New York Central build in all their required cities
	// before they build in any others.

	return errs
}
