package gameState

var companyInitCond = map[string]struct {
	tech3  bool
	sort   string
	tracks int
}{
	"New York, New Haven & Hartford": {tech3: false, sort: "00-00-01", tracks: 4},
	"Pennsylvania":                   {tech3: false, sort: "00-00-02", tracks: 16},
	"Boston & Maine":                 {tech3: false, sort: "00-00-03", tracks: 6},
	"Chesapeake & Ohio":              {tech3: false, sort: "00-00-04", tracks: 22},
	"New York Central":               {tech3: false, sort: "00-00-05", tracks: 20},
	"Baltimore & Ohio":               {tech3: false, sort: "00-00-06", tracks: 18},

	"New York, Chicago & Saint Louis": {tech3: true, sort: "00-00-07", tracks: 10},
	"Illinois Central":                {tech3: true, sort: "00-00-08", tracks: 12},
	"Erie":                            {tech3: true, sort: "00-00-09", tracks: 14},
	"Wabash":                          {tech3: true, sort: "00-00-10", tracks: 8},
}

func NewGame(playerNames []string) *Game {
	result := new(Game)

	result.GlobalState.TechLevel = 1
	result.GlobalState.UnminedCoal = []string{"G18", "H17", "I16", "J15", "K14"}
	result.GlobalState.OrphanStocks = make(map[string]int)

	result.Companies = make(map[string]*Company, len(companyInitCond))
	for name, start := range companyInitCond {
		result.Companies[name] = new(Company)
		result.Companies[name].Name = name
		result.Companies[name].HeldStock = 10
		result.Companies[name].Restricted = start.tech3
		result.Companies[name].PriceChange = start.sort
		result.Companies[name].UnbuiltTrack = start.tracks
		result.Companies[name].BuiltTrack = []string{}
	}

	// TODO: actually figure out how much cash each player is supposed to start with
	var startingCash int
	switch len(playerNames) {
	case 3, 4, 5, 6:
		startingCash = 250
	default:
		startingCash = 10
	}

	result.Players = make(map[string]*Player, len(playerNames))
	for _, name := range playerNames {
		result.Players[name] = new(Player)
		result.Players[name].Name = name
		result.Players[name].Cash = startingCash
		result.Players[name].Stocks = map[string]int{}
		result.Players[name].NetWorth = startingCash
	}

	result.beginMarketPhase()
	return result
}
