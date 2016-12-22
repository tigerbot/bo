package gameState

// The GlobalState struct holds all general board game state.
type GlobalState struct {
	Round int `json:"round"`
	Phase int `json:"phase"`

	TurnOrder  []string `json:"turn_order"`
	TurnNumber int      `json:"-"`
	Passes     int      `json:"-"`

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

// The Game struct holds all information for an active game.
type Game struct {
	GlobalState
	Companies map[string]*Company
	Players   map[string]*Player
}

// The MarketAction struct represents a single action that can be performed during the market
// phase. It represents the sale or purchase of stock from a single company.
type MarketAction struct {
	Company string `json:"company"`
	Count   int    `json:"count"`
	Price   int    `json:"price"`
}

// The MarketTurn struct represents all the possibilities for what can be done in a single market
// turn. The player can sell as many stocks as they have, but are only allowed to buy from one
// company during a single turn.
//
// TODO: players should be able to use a company's treasury to recover orphaned stock from the
// bank if they are the president. This action can only be done for one company per turn similar
// to buying stock, and cannot be done on the same turn the player uses to manage their own stock.
type MarketTurn struct {
	Sales    []MarketAction `json:"sales,omitempty"`
	Purchase *MarketAction  `json:"purchase,omitempty"`
}
