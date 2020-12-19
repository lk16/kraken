package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/lk16/noarray"
)

func unmarshalArray(bytes []byte, expectedLength int) ([]interface{}, error) {
	var slice []interface{}

	if err := json.Unmarshal(bytes, &slice); err != nil {
		return nil, fmt.Errorf("expeted JSON array")
	}

	if expectedLength != len(slice) {
		return nil, fmt.Errorf("expected JSON array with length %d", expectedLength)
	}

	return slice, nil
}

func unmarshalNumberasInt64(value interface{}) (int64, error) {
	float, ok := value.(float64)

	if !ok {
		return 0, errors.New("expected JSON number")
	}

	return int64(float), nil
}

func unmarshalStringasFloat64(value interface{}) (float64, error) {
	str, ok := value.(string)

	if !ok {
		return 0.0, errors.New("expected JSON string")
	}

	float, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return 0.0, fmt.Errorf("could not parse JSON string as float64: %w", err)
	}

	return float, nil
}

func (array *arrayModel) UnmarshalJSON(bytes []byte) error {
	var rawSlice []interface{}

	if err := json.Unmarshal(bytes, &rawSlice); err != nil {
		return err
	}

	if len(rawSlice) == 5 {
		// Usually the slice length is 4, however for BookUpdate values we CAN get slices with length 5.
		// In this case we get a separate JSON object at offset 2, other fields after that are as usual.
		// For ease of handling we just merge offset 2 into offset 1

		offset1, ok1 := rawSlice[1].(map[string]interface{})
		offset2, ok2 := rawSlice[2].(map[string]interface{})

		if !ok1 || !ok2 {
			return errors.New("expected JSON object at offsets 1 and 2")
		}

		for key, value := range offset2 {
			offset1[key] = value
		}

		// remove offset 2
		rawSlice = append(rawSlice[:2], rawSlice[3:]...)
	}

	if len(rawSlice) != 4 {
		return errors.New("expected JSON array with size 4")
	}

	var err error
	if array.ChannelID, err = unmarshalNumberasInt64(rawSlice[0]); err != nil {
		return fmt.Errorf("at position 0: %w", err)
	}

	array.Data = rawSlice[1]

	var ok bool
	if array.ChannelName, ok = rawSlice[2].(string); !ok {
		return errors.New("expected JSON string at offset 2")
	}

	if array.Pair, ok = rawSlice[3].(string); !ok {
		return errors.New("expected JSON string at offset 3")
	}

	return nil
}

func (ticker *Ticker) UnmarshalJSON(bytes []byte) error {

	rawSlice, err := unmarshalArray(bytes, 4)
	if err != nil {
		return err
	}

	if ticker.ChannelID, err = unmarshalNumberasInt64(rawSlice[0]); err != nil {
		return fmt.Errorf("at position 0: %w", err)
	}

	var ok bool
	if ticker.ChannelName, ok = rawSlice[2].(string); !ok {
		return errors.New("expected JSON string at offset 2")
	}

	if ticker.Pair, ok = rawSlice[3].(string); !ok {
		return errors.New("expected JSON string at offset 3")
	}

	return nil
}

func (tickerAskBid *TickerAskBid) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 3)
	if err != nil {
		return err
	}

	if tickerAskBid.Price, err = unmarshalStringasFloat64(rawSlice[0]); err != nil {
		return fmt.Errorf("at offset 0: %w", err)
	}

	if tickerAskBid.WholeLotVolume, err = unmarshalNumberasInt64(rawSlice[1]); err != nil {
		return fmt.Errorf("at offset 1: %w", err)
	}

	if tickerAskBid.LotVolume, err = unmarshalStringasFloat64(rawSlice[2]); err != nil {
		return fmt.Errorf("at offset 2: %w", err)
	}

	return nil
}

func (tickerClose *TickerClose) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 2)

	if err != nil {
		return err
	}

	if tickerClose.Price, err = unmarshalStringasFloat64(rawSlice[0]); err != nil {
		return fmt.Errorf("at offset 0: %w", err)
	}

	if tickerClose.LotVolume, err = unmarshalStringasFloat64(rawSlice[1]); err != nil {
		return fmt.Errorf("at offset 1: %w", err)
	}

	return nil
}

func (tickerTrades *TickerTrades) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 2)

	if err != nil {
		return err
	}

	if tickerTrades.Today, err = unmarshalNumberasInt64(rawSlice[0]); err != nil {
		return fmt.Errorf("at offset 0: %w", err)
	}

	if tickerTrades.Last24Hours, err = unmarshalNumberasInt64(rawSlice[1]); err != nil {
		return fmt.Errorf("at offset 1: %w", err)
	}

	return nil
}

func (tickerFloatStats *TickerFloatStats) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 2)

	if err != nil {
		return err
	}

	if tickerFloatStats.Today, err = unmarshalStringasFloat64(rawSlice[0]); err != nil {
		return fmt.Errorf("at offset 0: %w", err)
	}

	if tickerFloatStats.Last24Hours, err = unmarshalStringasFloat64(rawSlice[1]); err != nil {
		return fmt.Errorf("at offset 1: %w", err)
	}

	return nil
}

func (ohlcData *OHLCData) UnmarshalJSON(bytes []byte) error {
	return noarray.UnmarshalAsObject(bytes, ohlcData)
}

func (tradeData *TradeData) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 6)
	if err != nil {
		return err
	}

	if tradeData.Price, err = unmarshalStringasFloat64(rawSlice[0]); err != nil {
		return fmt.Errorf("at offset 0: %w", err)
	}

	if tradeData.Volume, err = unmarshalStringasFloat64(rawSlice[1]); err != nil {
		return fmt.Errorf("at offset 1: %w", err)
	}

	var tradeTime float64
	if tradeTime, err = unmarshalStringasFloat64(rawSlice[2]); err != nil {
		return fmt.Errorf("at offset 2: %w", err)
	}

	sec, dec := math.Modf(tradeTime)
	tradeData.Time = time.Unix(int64(sec), int64(dec*(1e9)))

	var ok bool
	if tradeData.Side, ok = rawSlice[3].(string); !ok {
		return errors.New("expected JSON string at offset 3")
	}

	if tradeData.OrderType, ok = rawSlice[4].(string); !ok {
		return errors.New("expected JSON string at offset 4")
	}

	if tradeData.Misc, ok = rawSlice[5].(string); !ok {
		return errors.New("expected JSON string at offset 5")
	}

	return nil
}

func (spreadData *SpreadData) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 5)

	if err != nil {
		return err
	}

	var timestamp float64

	floatPointers := []*float64{
		&spreadData.Ask, &spreadData.Bid, &timestamp,
		&spreadData.BidVolume, &spreadData.AskVolume,
	}

	for offset := range floatPointers {
		var value float64
		if value, err = unmarshalStringasFloat64(rawSlice[offset]); err != nil {
			return fmt.Errorf("at offset %d: %w", offset, err)
		}
		*floatPointers[offset] = value
	}

	sec, dec := math.Modf(timestamp)
	spreadData.Time = time.Unix(int64(sec), int64(dec*(1e9)))

	return nil
}

func (priceLevel *PriceLevel) UnmarshalJSON(bytes []byte) error {

	var rawSlice []interface{}
	err := json.Unmarshal(bytes, &rawSlice)

	if len(rawSlice) != 3 && len(rawSlice) != 4 {
		return errors.New("expected JSON array with length 3 or 4")
	}

	if err != nil {
		return err
	}

	var timestamp float64

	floatPointers := []*float64{&priceLevel.Price, &priceLevel.Volume, &timestamp}

	for offset := range floatPointers {
		var value float64
		if value, err = unmarshalStringasFloat64(rawSlice[offset]); err != nil {
			return fmt.Errorf("at offset %d: %w", offset, err)
		}
		*floatPointers[offset] = value
	}

	sec, dec := math.Modf(timestamp)
	priceLevel.Timestamp = time.Unix(int64(sec), int64(dec*(1e9)))

	return nil
}

func unmarshalArrayMessage(bytes []byte) (interface{}, error) {
	var array arrayModel

	if err := json.Unmarshal(bytes, &array); err != nil {
		return nil, fmt.Errorf("could not unmarshal into arrayModel: %w", err)
	}

	dataBytes, err := json.Marshal(array.Data)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JSON at offset 1: %w", err)
	}

	channelNamePrefix := strings.Split(array.ChannelName, "-")[0]

	switch channelNamePrefix {
	case "ticker":
		ticker := Ticker{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &ticker.Data); err != nil {
			return nil, fmt.Errorf("parsing Ticker data: %w", err)
		}
		return ticker, nil
	case "ohlc":
		ohlc := OHLC{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &ohlc.Data); err != nil {
			return nil, fmt.Errorf("parsing OHLC data: %w", err)
		}
		return ohlc, nil
	case "trade":
		trade := Trade{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &trade.Data); err != nil {
			return nil, fmt.Errorf("parsing trade data: %w", err)
		}
		return trade, nil
	case "spread":
		spread := Spread{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &spread.Data); err != nil {
			return nil, fmt.Errorf("parsing spread data: %w", err)
		}
		return spread, nil
	case "book":
		book := Book{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &book.Data); err != nil {
			return nil, fmt.Errorf("parsing book data: %w", err)
		}

		if book.Data.Asks == nil {
			bookUpdate := BookUpdate{
				ChannelID:   array.ChannelID,
				ChannelName: array.ChannelName,
				Pair:        array.Pair,
			}
			if err = json.Unmarshal(dataBytes, &bookUpdate.Data); err != nil {
				return nil, fmt.Errorf("parsing book update data: %w", err)
			}
			return bookUpdate, nil
		}

		return book, nil
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
