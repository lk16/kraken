package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookUpdate(t *testing.T) {

	type testCase struct {
		name         string
		book         Book
		update       BookUpdate
		expectedBook Book
	}

	testCases := []testCase{
		{
			"empty",
			Book{},
			BookUpdate{},
			Book{},
		},
		{
			"updateExisting",
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 2, Volume: 2},
				{Price: 4, Volume: 4},
			}}},
			BookUpdate{Data: BookUpdateData{Asks: []PriceLevel{
				{Price: 2, Volume: 99},
			}}},
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 2, Volume: 99},
				{Price: 4, Volume: 4},
			}}},
		},
		{
			"newLevel",
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 2, Volume: 2},
				{Price: 4, Volume: 4},
			}}},
			BookUpdate{Data: BookUpdateData{Asks: []PriceLevel{
				{Price: 3, Volume: 99},
			}}},
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 2, Volume: 2},
				{Price: 3, Volume: 99},
				{Price: 4, Volume: 4},
			}}},
		},
		{
			"deleteLevel",
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 2, Volume: 2},
				{Price: 3, Volume: 3},
				{Price: 4, Volume: 4},
			}}},
			BookUpdate{Data: BookUpdateData{Asks: []PriceLevel{
				{Price: 2, Volume: 0},
			}}},
			Book{Data: BookData{Asks: []PriceLevel{
				{Price: 3, Volume: 3},
				{Price: 4, Volume: 4},
			}}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.book.Update(testCase.update)
			assert.Equal(t, testCase.expectedBook, testCase.book)
		})
	}
}
