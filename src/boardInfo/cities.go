package boardInfo

import (
	"sort"
)

type City struct {
	Location string `json:"-"`
	Name     string `json:"name"`
	Revenue  [6]int `json:"revenue"`

	Exception string `json:"exception,omitempty"`
	Starting  string `json:"starting,omitempty"`
}

var cities = map[string]City{
	"A30": City{
		Name:    "Augusta",
		Revenue: [6]int{20, 20, 20, 20, 30, 40},
	},
	"B27": City{
		Name:    "Burlington",
		Revenue: [6]int{10, 20, 20, 20, 30, 30},
	},
	"C28": City{
		Name:    "Concord",
		Revenue: [6]int{20, 20, 20, 20, 20, 30},
	},
	"C30": City{
		Name:    "Portsmouth",
		Revenue: [6]int{20, 20, 20, 20, 20, 30},
	},
	"D19": City{
		Name:      "Buffalo",
		Revenue:   [6]int{20, 30, 30, 40, 50, 60},
		Exception: "New York Central",
		Starting:  "Erie",
	},
	"D21": City{
		Name:      "Syracuse",
		Revenue:   [6]int{10, 20, 20, 30, 30, 40},
		Exception: "New York Central",
	},
	"D23": City{
		Name:      "Utica",
		Revenue:   [6]int{10, 10, 10, 20, 20, 20},
		Exception: "New York Central",
	},
	"D25": City{
		Name:      "Albany",
		Revenue:   [6]int{30, 30, 40, 40, 40, 50},
		Exception: "New York Central",
		Starting:  "New York Central",
	},
	"D29": City{
		Name:     "Boston",
		Revenue:  [6]int{30, 30, 40, 40, 40, 50},
		Starting: "Boston & Maine",
	},
	"E4": City{
		Name:      "Chicago",
		Revenue:   [6]int{20, 30, 50, 70, 90, 100},
		Exception: "universal",
		Starting:  "New York, Chicago & Saint Louis",
	},
	"E12": City{
		Name:    "Detroit",
		Revenue: [6]int{20, 30, 40, 60, 80, 90},
	},
	"E28": City{
		Name:     "Hartford",
		Revenue:  [6]int{20, 20, 20, 30, 30, 30},
		Starting: "New York, New Haven & Hartford",
	},
	"E30": City{
		Name:    "Providence",
		Revenue: [6]int{20, 30, 30, 30, 30, 30},
	},
	"F13": City{
		Name:    "Cleveland",
		Revenue: [6]int{20, 30, 40, 50, 60, 60},
	},
	"F25": City{
		Name:      "New York",
		Revenue:   [6]int{30, 40, 50, 60, 70, 80},
		Exception: "New York Central",
	},
	"F27": City{
		Name:    "New Haven",
		Revenue: [6]int{20, 20, 30, 30, 30, 40},
	},
	"G2": City{
		Name:    "Springfield",
		Revenue: [6]int{10, 10, 20, 20, 20, 30},
	},
	"G8": City{
		Name:     "Fort Wayne",
		Revenue:  [6]int{10, 20, 20, 30, 40, 50},
		Starting: "Wabash",
	},
	"G16": City{
		Name:      "Pittsburgh",
		Revenue:   [6]int{20, 30, 40, 60, 70, 80},
		Exception: "Pennsylvania",
	},
	"G20": City{
		Name:      "Harrisburg",
		Revenue:   [6]int{10, 10, 20, 20, 20, 20},
		Exception: "Pennsylvania",
	},
	"G24": City{
		Name:      "Philadelphia",
		Revenue:   [6]int{30, 40, 40, 40, 50, 60},
		Exception: "Pennsylvania",
		Starting:  "Pennsylvania",
	},
	"H7": City{
		Name:    "Indianapolis",
		Revenue: [6]int{20, 30, 30, 40, 50, 60},
	},
	"H15": City{
		Name:    "Wheeling",
		Revenue: [6]int{20, 20, 30, 40, 50, 60},
	},
	"H23": City{
		Name:     "Baltimore",
		Revenue:  [6]int{20, 30, 30, 40, 40, 50},
		Starting: "Baltimore & Ohio",
	},
	"I0": City{
		Name:     "Saint Louis",
		Revenue:  [6]int{30, 40, 50, 60, 70, 90},
		Starting: "Illinois Central",
	},
	"I10": City{
		Name:    "Cincinnati",
		Revenue: [6]int{30, 40, 50, 50, 60, 70},
	},
	"I22": City{
		Name:    "Washington",
		Revenue: [6]int{20, 20, 30, 30, 30, 30},
	},
	"I24": City{
		Name:    "Dover",
		Revenue: [6]int{10, 10, 10, 20, 20, 20},
	},
	"J7": City{
		Name:    "Louisville",
		Revenue: [6]int{20, 30, 30, 40, 40, 50},
	},
	"J13": City{
		Name:    "Huntington",
		Revenue: [6]int{10, 10, 20, 30, 30, 40},
	},
	"J21": City{
		Name:     "Richmond",
		Revenue:  [6]int{30, 30, 20, 20, 20, 30},
		Starting: "Chesapeake & Ohio",
	},
	"K2": City{
		Name:    "Cairo",
		Revenue: [6]int{10, 20, 20, 20, 20, 20},
	},
	"K10": City{
		Name:    "Lexington",
		Revenue: [6]int{10, 20, 20, 30, 30, 30},
	},
	"K16": City{
		Name:    "Roanoke",
		Revenue: [6]int{20, 20, 20, 20, 20, 20},
	},
	"K22": City{
		Name:    "Norfolk",
		Revenue: [6]int{20, 20, 30, 30, 30, 40},
	},
}

// startingLocations is used as a quick look-up for the StartingLocation function
var startingLocations map[string]string

func init() {
	startingLocations = make(map[string]string, 10)
	for coord, val := range cities {
		// In order to contain the location in the struct and not require having two duplicate
		// values in the map definition we assign it here. Golang does not allow us to assign
		// to members directory in the map, so we have to modify the value and overwrite the
		// previous map entry with the modified one.
		val.Location = coord
		cities[coord] = val

		if val.Starting != "" {
			startingLocations[val.Starting] = coord
		}
	}
}

// StartingLocation looks the map coordinate for the specified company's starting location.
func StartingLocation(company string) string {
	return startingLocations[company]
}

// Cities returns a slice of all the cities that coincide with the provided map coordinates
func Cities(mapCoords ...string) []City {
	result := make([]City, 0, len(mapCoords)/2)

	for _, coord := range mapCoords {
		if val, ok := cities[coord]; ok {
			result = append(result, val)
		}
	}

	return result
}

func SortCities(cityList []City, techLvl int) {
	sort.Sort(citySorter{cityList, techLvl})
}

// citySorter implements sort.Interface and sorts the list of cities based on revenue for the
// given tech level. Cities with the highest revenues are placed in the beginning of the list.
type citySorter struct {
	cityList []City
	techLvl  int
}

func (s citySorter) Len() int {
	return len(s.cityList)
}
func (s citySorter) Less(i, j int) bool {
	return s.cityList[i].Revenue[s.techLvl] < s.cityList[j].Revenue[s.techLvl]
}
func (s citySorter) Swap(i, j int) {
	s.cityList[i], s.cityList[j] = s.cityList[j], s.cityList[i]
}
