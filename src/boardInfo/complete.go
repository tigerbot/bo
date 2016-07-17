// Package boardInfo contains all of the static information for the board that should not change
// between games or at any point during a game. Because go does not allow anything but primitives
// to be constants this is a separate package largely to make sure this information is never
// accidently modified.
package boardInfo

import (
	"encoding/json"
)

type hex struct {
	BuildCost int   `json:"build_cost"`
	City      *City `json:"city"`
}

var completeMap map[string]hex

func init() {
	completeMap = make(map[string]hex, len(buildCosts))

	for coord, cost := range buildCosts {
		var curCity *City
		if val, ok := cities[coord]; ok {
			curCity = &val
		}
		completeMap[coord] = hex{
			BuildCost: cost,
			City:      curCity,
		}
	}
}

func JsonMap() ([]byte, error) {
	return json.Marshal(completeMap)
}
