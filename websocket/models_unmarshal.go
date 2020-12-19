package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lk16/noarray"
)

func (array *arrayModel) UnmarshalJSON(bytes []byte) error {
	var rawSlice []interface{}

	if err := json.Unmarshal(bytes, &rawSlice); err != nil {
		return err
	}

	if len(rawSlice) < 2 {
		return errors.New("Expected JSON array of at least 3")
	}

	channelName, ok := rawSlice[len(rawSlice)-2].(string)
	if !ok {
		return fmt.Errorf("expected string at offset %d", len(rawSlice)-2)
	}

	array.ChannelName = channelName
	return nil
}

func (ticker *Ticker) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, ticker)
}

func (tickerAskBid *TickerAskBid) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, tickerAskBid)
}

func (tickerClose *TickerClose) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, tickerClose)
}

func (tickerTrades *TickerTrades) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, tickerTrades)
}

func (tickerFloatStats *TickerFloatStats) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, tickerFloatStats)
}

func (ohlcData *OHLCData) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, ohlcData)
}

func (ohlc *OHLC) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, ohlc)
}

func (trade *Trade) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, trade)
}

func (spread *Spread) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, spread)
}

func (tradeData *TradeData) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, tradeData)
}

func (spreadData *SpreadData) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, spreadData)
}

func (book *Book) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, book)
}

func (priceLevel *PriceLevel) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, priceLevel)
}

func (bookUpdate *BookUpdate) UnmarshalJSON(bytes []byte) error {
	if err := noarray.UnmarshalAsObject(bytes, bookUpdate); err == nil {
		return nil
	}

	// Asks and bids can arrive in the same update.
	// When this happens they come in a different object (for
	// some reason) at offsets 1 and 2. We merge them here
	// and proceed as-if they came in one JSON object.

	type bookUpdateFiveItems struct {
		ChannelID   int64          `json:"0"`
		Asks        BookUpdateData `json:"1"`
		Bids        BookUpdateData `json:"2"`
		ChannelName string         `json:"3"`
		Pair        string         `json:"4"`
	}

	var bookUpdateAsksAndBids bookUpdateFiveItems
	if err := noarray.UnmarshalAsObject(bytes, &bookUpdateAsksAndBids); err != nil {
		return err
	}

	bookUpdate.ChannelID = bookUpdateAsksAndBids.ChannelID
	bookUpdate.Data.Asks = bookUpdateAsksAndBids.Asks.Asks
	bookUpdate.Data.Bids = bookUpdateAsksAndBids.Bids.Bids
	bookUpdate.ChannelName = bookUpdateAsksAndBids.ChannelName
	bookUpdate.Pair = bookUpdateAsksAndBids.Pair

	return nil
}

func unmarshalArrayMessage(bytes []byte) (interface{}, error) {
	var array arrayModel

	if err := json.Unmarshal(bytes, &array); err != nil {
		return nil, fmt.Errorf("could not unmarshal into arrayModel: %w", err)
	}

	var err error
	channelNamePrefix := strings.Split(array.ChannelName, "-")[0]

	switch channelNamePrefix {
	case "ticker":
		var ticker Ticker
		if err = json.Unmarshal(bytes, &ticker); err != nil {
			return nil, fmt.Errorf("parsing Ticker: %w", err)
		}
		return ticker, nil
	case "ohlc":
		var ohlc OHLC
		if err = json.Unmarshal(bytes, &ohlc); err != nil {
			return nil, fmt.Errorf("parsing OHLC: %w", err)
		}
		return ohlc, nil
	case "trade":
		var trade Trade
		if err = json.Unmarshal(bytes, &trade); err != nil {
			return nil, fmt.Errorf("parsing Trade: %w", err)
		}
		return trade, nil
	case "spread":
		var spread Spread
		if err = json.Unmarshal(bytes, &spread); err != nil {
			return nil, fmt.Errorf("parsing Spread: %w", err)
		}
		return spread, nil
	case "book":
		var book Book
		err = json.Unmarshal(bytes, &book)
		if err == nil && book.Data.Asks != nil {
			return book, nil
		}

		var bookUpdate BookUpdate
		if err = json.Unmarshal(bytes, &bookUpdate); err != nil {
			return nil, fmt.Errorf("parsing Book/BookUpdate: %w", err)
		}
		return bookUpdate, nil

	default:
		return nil, fmt.Errorf("unknown channel name prefix %s", channelNamePrefix)
	}
}

func unmarshalReceivedMessage(bytes []byte) (interface{}, error) {

	var event event

	if err := json.Unmarshal(bytes, &event); err != nil {
		// This probably means the message is not a JSON object.
		// All other kraken models are JSON arrays, so we try those.
		// The case of broken JSON also ends up here.
		return unmarshalArrayMessage(bytes)
	}

	if event.Status == "error" {
		var model Error
		if err := json.Unmarshal(bytes, &model); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	}

	switch event.Event {
	case "heartbeat":
		var model HeartBeat
		if err := json.Unmarshal(bytes, &model); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	case "pong":
		var model Pong
		if err := json.Unmarshal(bytes, &model); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	case "subscriptionStatus":
		var model SubscriptionStatus
		if err := json.Unmarshal(bytes, &model); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	case "systemStatus":
		var model SystemStatus
		if err := json.Unmarshal(bytes, &model); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	default:
		return nil, fmt.Errorf("unknown model %s", event.Event)
	}
}
