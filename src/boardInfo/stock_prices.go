package boardInfo

import (
	"sort"
)

var stockPrices = []int{
	34,
	37,
	41,
	45,
	50,
	55,
	60,
	66,
	74,
	82,
	91,
	100,
	110,
	121,
	133,
	148,
	160,
	176,
	194,
	213,
	234,
	257,
	282,
	310,
	341,
	375,
}

func NextStockPrice(curPrice int) int {
	ind := sort.SearchInts(stockPrices, curPrice)
	if stockPrices[ind] == curPrice && ind+1 < len(stockPrices) {
		return stockPrices[ind+1]
	} else if ind < len(stockPrices) {
		return stockPrices[ind]
	} else {
		return stockPrices[len(stockPrices)-1]
	}
}

func PrevStockPrice(curPrice int) int {
	ind := sort.SearchInts(stockPrices, curPrice)
	if ind-1 >= 0 {
		return stockPrices[ind-1]
	} else {
		return stockPrices[0]
	}
}

func StartingStockPrices(techLevel int) [3]int {
	return [3]int{
		stockPrices[4+techLevel],
		stockPrices[5+techLevel],
		stockPrices[6+techLevel],
	}
}
