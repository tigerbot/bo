package gameState

// The GlobalState struct holds all general board game state.
type GlobalState struct {
	Round int `json:"round"`
	Phase int `json:"phase"`

	turnOrder []string `json:"-"`
	turn      int      `json:"-"`
	passes    int      `json:"-"`

	TrainsBought int            `json:"trains_bought"`
	TechLevel    int            `json:"tech_level"`
	UnminedCoal  []string       `json:"unmined_coal"`
	OrphanStocks map[string]int `json:"orphan_stocks"`
}

// The Company struct holds all of the information relevant to a single company.
type Company struct {
	Name       string `json:"-"`
	Restricted bool   `json:"restricted"`
	President  string `json:"president"`

	StockPrice  int    `json:"stock_price"`
	PriceChange string `json:"price_changed"`
	HeldStock   int    `json:"held_stock"`
	NetIncome   int    `json:"net_income"`
	Treasury    int    `json:"treasury"`

	CoalMined    int      `json:"coal_mined"`
	UnbuiltTrack int      `json:"unbuilt_track"`
	BuiltTrack   []string `json:"built_track"`
	Equipment    [6]int   `json:"equipment"`
}

// The Player struct keeps track of a single players liquid assets and stock.
type Player struct {
	Name     string         `json:"-"`
	Cash     int            `json:"cash"`
	NetWorth int            `json:"net_worth"`
	Stocks   map[string]int `json:"stocks"`
}

type Game struct {
	GlobalState
	Companies map[string]*Company
	Players   map[string]*Player
}
