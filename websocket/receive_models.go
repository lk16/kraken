package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
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

type event struct {
	Event string `json:"event"`
}

type SystemStatus struct {
	Event        string      `json:"event"`
	ConnectionID json.Number `json:"connectionID"`
	Status       string      `json:"status"`
	Version      string      `json:"version"`
}

type SubscriptionDetails struct {
	Depth        int    `json:"depth,omitempty"`
	Interval     int    `json:"interval,omitempty"`
	MaxRateCount int    `json:"maxratecount,omitempty"`
	Name         string `json:"name"`
	Token        string `json:"token,omitempty"`
}

type SubscriptionStatus struct {
	Event        string              `json:"event"`
	ChannelID    int                 `json:"channelID"`
	ChannelName  string              `json:"channelName"`
	ReqID        int                 `json:"reqid"`
	Pair         string              `json:"pair"`
	Status       string              `json:"status"`
	Subscription SubscriptionDetails `json:"subscription"`
	ErrorMessage string              `json:"errorMessage"`
}

type Pong struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid"`
}

type HeartBeat struct {
	Event string `json:"event"`
}

type Ticker struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        TickerData
}

type arrayModel struct {
	ChannelID   int64
	Data        interface{}
	ChannelName string
	Pair        string
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

type TickerData struct {
	Ask                   TickerAskBid     `json:"a"`
	Bid                   TickerAskBid     `json:"b"`
	Close                 TickerClose      `json:"c"`
	Trades                TickerTrades     `json:"t"`
	Volume                TickerFloatStats `json:"v"`
	VolumeWeightedAverage TickerFloatStats `json:"p"`
	Low                   TickerFloatStats `json:"l"`
	High                  TickerFloatStats `json:"h"`
	Open                  TickerFloatStats `json:"o"`
}

type TickerAskBid struct {
	Price          float64
	WholeLotVolume int64
	LotVolume      float64
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

type TickerClose struct {
	Price     float64
	LotVolume float64
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

type TickerTrades struct {
	Today       int64
	Last24Hours int64
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

type TickerFloatStats struct {
	Today       float64
	Last24Hours float64
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

type OHLC struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        OHLCData
}

type OHLCData struct {
	Time                time.Time
	EndTime             time.Time
	Open                float64
	High                float64
	Low                 float64
	Close               float64
	VolumeWeightedPrice float64
	Volume              float64
	Count               int64
}

func (ohlcData *OHLCData) UnmarshalJSON(bytes []byte) error {
	rawSlice, err := unmarshalArray(bytes, 9)

	if err != nil {
		return err
	}

	var currentTime, endTime float64

	floatPointers := []*float64{
		&currentTime, &endTime, &ohlcData.Open, &ohlcData.High, &ohlcData.Low,
		&ohlcData.Close, &ohlcData.VolumeWeightedPrice, &ohlcData.Volume,
	}

	for offset := range floatPointers {
		var value float64
		if value, err = unmarshalStringasFloat64(rawSlice[offset]); err != nil {
			return fmt.Errorf("at offset %d: %w", offset, err)
		}
		*floatPointers[offset] = value
	}

	sec, dec := math.Modf(currentTime)
	ohlcData.Time = time.Unix(int64(sec), int64(dec*(1e9)))

	sec, dec = math.Modf(endTime)
	ohlcData.EndTime = time.Unix(int64(sec), int64(dec*(1e9)))

	if ohlcData.Count, err = unmarshalNumberasInt64(rawSlice[8]); err != nil {
		return fmt.Errorf("at offset 8: %w", err)
	}

	return nil
}

type Trade struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        []TradeData
}

type TradeData struct {
	Price     float64
	Volume    float64
	Time      time.Time
	Side      string
	OrderType string
	Misc      string
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

type Spread struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        SpreadData
}

type SpreadData struct {
	Ask       float64
	Bid       float64
	Time      time.Time
	BidVolume float64
	AskVolume float64
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

type Book struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        BookData
}

func (book *Book) Print() {
	fmt.Printf("Asks:\n")
	for _, ask := range book.Data.Asks {
		fmt.Printf("%11.5f %11.5f\n", ask.Price, ask.Volume)
	}

	fmt.Printf("Bids:\n")
	for _, bid := range book.Data.Bids {
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
			if removeLevel {
				log.Printf("WARNING: attempting to remove non-existent level")
			} else {
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

type BookData struct {
	Asks []PriceLevel `json:"as"`
	Bids []PriceLevel `json:"bs"`
}

type BookUpdate struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        BookUpdateData
}

type BookUpdateData struct {
	Asks []PriceLevel `json:"a"`
	Bids []PriceLevel `json:"b"`
}

type PriceLevel struct {
	Price     float64
	Volume    float64
	Timestamp time.Time
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
		return nil, fmt.Errorf("re-marshaling JSON at position 1: %w", err)
	}

	channelNamePrefix := strings.Split(array.ChannelName, "-")[0]

	switch channelNamePrefix {
	case "ticker":
		ticker := &Ticker{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &ticker.Data); err != nil {
			return nil, fmt.Errorf("parsing Ticker data: %w", err)
		}
		return ticker, nil
	case "ohlc":
		ohlc := &OHLC{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &ohlc.Data); err != nil {
			return nil, fmt.Errorf("parsing OHLC data: %w", err)
		}
		return ohlc, nil
	case "trade":
		trade := &Trade{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &trade.Data); err != nil {
			return nil, fmt.Errorf("parsing trade data: %w", err)
		}
		return trade, nil
	case "spread":
		spread := &Spread{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &spread.Data); err != nil {
			return nil, fmt.Errorf("parsing spread data: %w", err)
		}
		return spread, nil
	case "book":
		book := &Book{
			ChannelID:   array.ChannelID,
			ChannelName: array.ChannelName,
			Pair:        array.Pair,
		}
		if err = json.Unmarshal(dataBytes, &book.Data); err != nil {
			return nil, fmt.Errorf("parsing book data: %w", err)
		}

		if book.Data.Asks == nil {
			bookUpdate := &BookUpdate{
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

	modelOrError := func(model interface{}, err error) (interface{}, error) {
		if err != nil {
			return nil, fmt.Errorf("parsing %s failed: %w", event.Event, err)
		}
		return model, nil
	}

	if err := json.Unmarshal(bytes, &event); err != nil {
		// This probably means the message is not a JSON object.
		// All other kraken models are JSON arrays, so we try those.
		// The case of broken JSON also ends up here.
		model, err := unmarshalArrayMessage(bytes)

		if err != nil {
			// TODO: remove
			log.Printf("WARNING: error while parsing message %s", string(bytes))
		}

		return model, err

	}

	switch event.Event {
	case "heartbeat":
		var model HeartBeat
		return modelOrError(&model, json.Unmarshal(bytes, &model))
	case "pong":
		var model Pong
		return modelOrError(&model, json.Unmarshal(bytes, &model))
	case "subscriptionStatus":
		var model SubscriptionStatus
		return modelOrError(&model, json.Unmarshal(bytes, &model))
	case "systemStatus":
		var model SystemStatus
		return modelOrError(&model, json.Unmarshal(bytes, &model))
	default:
		return modelOrError(nil, errors.New("unknown model"))
	}
}
