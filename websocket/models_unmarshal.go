package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func getMessageType(bytes []byte) (string, error) {
	var event event
	if err := json.Unmarshal(bytes, &event); err != nil {

		var array arrayModel
		if err := json.Unmarshal(bytes, &array); err != nil {
			return "", fmt.Errorf("could not unmarshal into arrayModel: %w", err)
		}

		return strings.Split(array.ChannelName, "-")[0], nil
	}

	if event.Status == "error" {
		return "error", nil
	}

	return event.Event, nil
}

func unmarshalReceivedMessage(bytes []byte) (interface{}, error) {

	messageType, err := getMessageType(bytes)

	if err != nil {
		return nil, err
	}

	if messageType == "book" {
		var book Book
		if err := json.Unmarshal(bytes, &book); err == nil && book.Data.Asks != nil {
			return book, nil
		}

		var bookUpdate BookUpdate
		if err = json.Unmarshal(bytes, &bookUpdate); err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", messageType, err)
		}
		return bookUpdate, nil
	}

	targetMap := map[string]interface{}{
		"error":              &Error{},
		"heartbeat":          &HeartBeat{},
		"pong":               &Pong{},
		"subscriptionStatus": &SubscriptionStatus{},
		"systemStatus":       &SystemStatus{},
		"ticker":             &Ticker{},
		"ohlc":               &OHLC{},
		"trade":              &Trade{},
		"spread":             &Spread{},
	}

	target, ok := targetMap[messageType]

	if !ok {
		return nil, fmt.Errorf("unknown message type %s", messageType)
	}

	if err := json.Unmarshal(bytes, &target); err != nil {
		return nil, fmt.Errorf("parsing %s failed: %w", messageType, err)
	}

	return reflect.Indirect(reflect.ValueOf(target)).Interface(), nil

}
