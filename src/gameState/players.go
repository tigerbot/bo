package gameState

type Player struct {
	Cash     int            `json:"cash"`
	NetWorth int            `json:"net_worth"`
	Stocks   map[string]int `json:"stocks"`
}
