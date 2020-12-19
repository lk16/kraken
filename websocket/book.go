package websocket

import (
	"fmt"
	"sort"
)

func (book *Book) PrintTop(n int) {
	fmt.Printf("Asks:\n")
	for _, ask := range book.Data.Asks[:n] {
		fmt.Printf("%11.5f %11.5f\n", ask.Price, ask.Volume)
	}

	fmt.Printf("Bids:\n")
	for _, bid := range book.Data.Bids[:n] {
		fmt.Printf("%11.5f %11.5f\n", bid.Price, bid.Volume)
	}
}

func updateSide(side []PriceLevel, updates []PriceLevel) []PriceLevel {
	for _, update := range updates {
		price := update.Price
		removeLevel := update.Volume == 0.0

		foundIndex := -1

		for index, level := range side {
			if level.Price == price {
				foundIndex = index
				break
			}
		}

		if foundIndex != -1 {
			if removeLevel {
				// swap with last
				side[len(side)-1], side[foundIndex] = side[foundIndex], side[len(side)-1]

				// remove last item
				side = side[:len(side)-1]

			} else {
				// update level
				side[foundIndex] = update
			}
		} else {
			if !removeLevel {
				// add level
				side = append(side, update)
			}
		}
	}
	return side
}

func (book *Book) Update(update BookUpdate) {
	book.Data.Asks = updateSide(book.Data.Asks, update.Data.Asks)
	book.Data.Bids = updateSide(book.Data.Bids, update.Data.Bids)

	sort.Slice(book.Data.Asks, func(i, j int) bool {
		return book.Data.Asks[i].Price < book.Data.Asks[j].Price
	})

	sort.Slice(book.Data.Bids, func(i, j int) bool {
		return book.Data.Bids[i].Price > book.Data.Bids[j].Price
	})
}
