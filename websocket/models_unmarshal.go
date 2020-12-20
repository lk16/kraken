package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (array *arrayModel) UnmarshalJSON(bytes []byte) error {
	var rawMessages []json.RawMessage

	if err := json.Unmarshal(bytes, &rawMessages); err != nil {
		return err
	}

	if len(rawMessages) < 2 {
		return errors.New("Expected JSON array of at least 3")
	}

	return json.Unmarshal(rawMessages[len(rawMessages)-2], &array.ChannelName)
}

func (ticker *Ticker) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&ticker.ChannelID,
		&ticker.Data,
		&ticker.ChannelName,
		&ticker.Pair,
	}
	return json.Unmarshal(bytes, &slice)
}

func (tickerAskBid *TickerAskBid) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&tickerAskBid.Price,
		&tickerAskBid.WholeLotVolume,
		&tickerAskBid.LotVolume,
	}
	return json.Unmarshal(bytes, &slice)
}

func (tickerClose *TickerClose) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&tickerClose.Price,
		&tickerClose.LotVolume,
	}
	return json.Unmarshal(bytes, &slice)
}

func (tickerTrades *TickerTrades) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&tickerTrades.Today,
		&tickerTrades.Last24Hours,
	}
	return json.Unmarshal(bytes, &slice)
}

func (tickerFloatStats *TickerFloatStats) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&tickerFloatStats.Today,
		&tickerFloatStats.Last24Hours,
	}
	return json.Unmarshal(bytes, &slice)
}

func (ohlcData *OHLCData) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&ohlcData.Time,
		&ohlcData.EndTime,
		&ohlcData.Open,
		&ohlcData.High,
		&ohlcData.Low,
		&ohlcData.Close,
		&ohlcData.VolumeWeightedPrice,
		&ohlcData.Volume,
		&ohlcData.Count,
	}
	return json.Unmarshal(bytes, &slice)
}

func (ohlc *OHLC) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&ohlc.ChannelID,
		&ohlc.Data,
		&ohlc.ChannelName,
		&ohlc.Pair,
	}
	return json.Unmarshal(bytes, &slice)
}

func (trade *Trade) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&trade.ChannelID,
		&trade.Data,
		&trade.ChannelName,
		&trade.Pair,
	}
	return json.Unmarshal(bytes, &slice)
}

func (spread *Spread) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&spread.ChannelID,
		&spread.Data,
		&spread.ChannelName,
		&spread.Pair,
	}
	return json.Unmarshal(bytes, &slice)
}

func (tradeData *TradeData) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&tradeData.Price,
		&tradeData.Volume,
		&tradeData.Time,
		&tradeData.Side,
		&tradeData.OrderType,
		&tradeData.Misc,
	}
	return json.Unmarshal(bytes, &slice)
}

func (spreadData *SpreadData) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&spreadData.Ask,
		&spreadData.Bid,
		&spreadData.Time,
		&spreadData.BidVolume,
		&spreadData.AskVolume,
	}
	return json.Unmarshal(bytes, &slice)
}

func (book *Book) UnmarshalJSON(bytes []byte) error {
	slice := []interface{}{
		&book.ChannelID,
		&book.Data,
		&book.ChannelName,
		&book.Pair,
	}
	return json.Unmarshal(bytes, &slice)
}

func (priceLevel *PriceLevel) UnmarshalJSON(bytes []byte) error {

	slice := []interface{}{
		&priceLevel.Price,
		&priceLevel.Volume,
		&priceLevel.Timestamp,
	}
	return json.Unmarshal(bytes, &slice)
}

func (bookUpdate *BookUpdate) UnmarshalJSON(bytes []byte) error {

	slice := []interface{}{
		&bookUpdate.ChannelID,
		&bookUpdate.Data,
		&bookUpdate.ChannelName,
		&bookUpdate.Pair,
	}

	if err := json.Unmarshal(bytes, &slice); err == nil {
		return nil
	}

	// Asks and bids can arrive in the same update.
	// When this happens they come in a separate object (for
	// some reason) at offsets 1 and 2. We merge them here
	// and proceed as-if they came in one JSON object.
	var separateBids BookUpdateData

	slice = []interface{}{
		&bookUpdate.ChannelID,
		&bookUpdate.Data,
		&separateBids,
		&bookUpdate.ChannelName,
		&bookUpdate.Pair,
	}

	if err := json.Unmarshal(bytes, &slice); err != nil {
		return err
	}

	bookUpdate.Data.Bids = separateBids.Bids
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
