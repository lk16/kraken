package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

func (ticker *Ticker) UnmarshalJSON(bytes []byte) error {

	rawSlice, err := unmarshalArray(bytes, 4)
	if err != nil {
		return err
	}

	if ticker.ChannelID, err = unmarshalNumberasInt64(rawSlice[0]); err != nil {
		return fmt.Errorf("at position 0: %w", err)
	}

	tickerDataBytes, err := json.Marshal(rawSlice[1])
	if err != nil {
		return fmt.Errorf("at position 1: %w", err)
	}

	if err = json.Unmarshal(tickerDataBytes, &ticker.Data); err != nil {
		return fmt.Errorf("parsing data object at offset 1 failed: %w", err)
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

func unmarshalArrayMessage(bytes []byte) (interface{}, error) {
	rawSlice, err := unmarshalArray(bytes, 4)
	if err != nil {
		return nil, err
	}

	kind, ok := rawSlice[2].(string)
	if !ok {
		return nil, fmt.Errorf("expected JSON string at offset 2")
	}

	if kind == "ticker" {
		var ticker Ticker
		var err error
		if err = json.Unmarshal(bytes, &ticker); err == nil {
			return &ticker, nil
		}
		return nil, err
	}

	return nil, fmt.Errorf("could not recognize JSON message with kind %s", kind)
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
		return unmarshalArrayMessage(bytes)
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
