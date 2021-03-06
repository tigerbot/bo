// Package boardInfo contains all of the static information for the board that should not change
// between games or at any point during a game. Because go does not allow anything but primitives
// to be constants this is a separate package largely to make sure this information is never
// accidentally modified.
package boardInfo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type hex struct {
	BuildCost int   `json:"build_cost"`
	City      *City `json:"city"`
	Coal      bool  `json:"coal"`
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

	for _, coord := range StartingCoal() {
		curHex := completeMap[coord]
		curHex.Coal = true
		completeMap[coord] = curHex
	}
}

func JsonMap() ([]byte, error) {
	return json.Marshal(completeMap)
}

func StartingCoal() []string {
	return []string{"G18", "H17", "I16", "J15", "K14"}
}

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

// TilesAdjacent checks if two map coordinates are next to each other.
func TilesAdjacent(coordA, coordB string) bool {
	// Handle the exceptions created by the impassable borders.
	if (coordA == "I22" && coordB == "I24") || (coordB == "I22" && coordA == "I24") {
		return false
	}
	if (coordA == "E12" && coordB == "F13") || (coordB == "E12" && coordA == "F13") {
		return false
	}

	// If the hexes aren't in the same or adjacent row then they cannot be adjacent. We have to
	// cast to int before the subtraction since the byte is an unsigned value.
	rowDiff := absInt(int(coordA[0]) - int(coordB[0]))
	if rowDiff > 1 {
		return false
	}

	var colDiff int
	if colA, err := strconv.Atoi(coordA[1:]); err != nil {
		fmt.Printf("failed to get column from hex coordinate %q: %v\n", coordA, err)
		return false
	} else if colB, err := strconv.Atoi(coordB[1:]); err != nil {
		fmt.Printf("failed to get column from hex coordinate %q: %v\n", coordB, err)
		return false
	} else {
		colDiff = absInt(int(colA - colB))
	}

	// Because of the way we handle the column indices (only using even or odd numbers within a
	// row) the colDiff will be 2 for adjacent tiles in the same row and 1 for adjacent tiles
	// in adjacent rows, so the sum should always be 2 ()
	return rowDiff+colDiff == 2
}

// tileConnected checks if a map coordinate has any adjacent tiles within a list of coordinates.
func tileConnected(coord string, blob []string) bool {
	for _, value := range blob {
		if TilesAdjacent(coord, value) {
			return true
		}
	}
	return false
}

// TilesContiguous checks if there would be an unbroken path from each coordinate in update to
// the coordinates in existent if all tiles in both lists were colored. No checking is done on
// the existent coordinates.
//
// This is a recursive function. If in any given call, there are coordinates that were not adjacent
// to any in the existent list there is no point in checking those relations again, so we create
// a new existent list with the updates that were established as contiguous and check the ones
// that have not yet been confirm against them.
func TilesContiguous(existent, update []string) bool {
	connected := make([]string, 0, len(update))
	isolated := make([]string, 0, len(update))
	for _, coord := range update {
		if tileConnected(coord, existent) {
			connected = append(connected, coord)
		} else {
			isolated = append(isolated, coord)
		}
	}

	if len(isolated) == 0 {
		// All tiles have been confirmed as connected, return true
		return true
	} else if len(connected) == 0 {
		// No addition tiles were establish as connected so there's nothing left to check
		return false
	}
	return TilesContiguous(connected, isolated)
}

// TechLevel converts the number of trains that have been bought during the game into the
// tech level. The conversion is rather simple, and this is its own function just to make
// sure that all the places that need to determine the tech level are consistent.
func TechLevel(trainsBought int) int {
	if trainsBought <= 0 {
		return 1
	}
	// We enter the tech level only after the first train of that level has been purchased
	return ((trainsBought - 1) / 5) + 1
}
