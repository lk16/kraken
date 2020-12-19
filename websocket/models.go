package websocket

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

type UnixTime time.Time

func (unixTime *UnixTime) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}

	var unixTimeFloat float64
	var err error
	if unixTimeFloat, err = strconv.ParseFloat(str, 64); err != nil {
		return fmt.Errorf("could not parse JSON string as float64: %w", err)
	}

	sec, dec := math.Modf(unixTimeFloat)
	*unixTime = UnixTime(time.Unix(int64(sec), int64(dec*(1e9))))
	return nil
}

type Float64String float64

func (float64string *Float64String) UnmarshalJSON(bytes []byte) error {
	var number json.Number
	if err := json.Unmarshal(bytes, &number); err != nil {
		return err
	}

	float, err := number.Float64()
	if err != nil {
		return err
	}

	*float64string = Float64String(float)
	return nil
}

type Int64String int64

func (int64String *Int64String) UnmarshalJSON(bytes []byte) error {
	var number json.Number
	if err := json.Unmarshal(bytes, &number); err != nil {
		return err
	}

	anInt64, err := number.Int64()
	if err != nil {
		return err
	}

	*int64String = Int64String(anInt64)
	return nil
}

type Ping struct {
	Event string `json:"event"`
	ReqID int    `json:"reqid,omitempty"`
}

type Subscription struct {
	Name        string `json:"name"`
	Depth       int    `json:"depth,omitempty"`
	Interval    int    `json:"interval,omitempty"`
	RateCounter bool   `json:"ratecounter,omitempty"`
	Snapshot    bool   `json:"snapshot,omitempty"`
	Token       string `json:"token,omitempty"`
}

type Subscribe struct {
	Event        string       `json:"event"`
	ReqID        int          `json:"reqid,omitempty"`
	Pair         []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

type Unsubscribe struct {
	Event        string       `json:"event"`
	ReqID        int          `json:"reqid"`
	Pair         []string     `json:"pair"`
	Subscription Subscription `json:"subscription"`
}

type event struct {
	Event  string `json:"event"`
	Status string `json:"status"`
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
	ChannelID   Int64String `json:"0"`
	Data        TickerData  `json:"1"`
	ChannelName string      `json:"2"`
	Pair        string      `json:"3"`
}

type arrayModel struct {
	ChannelName string
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
	Price          Float64String `json:"0"`
	WholeLotVolume Int64String   `json:"1"`
	LotVolume      Float64String `json:"2"`
}

type TickerTrades struct {
	Today       Int64String `json:"0"`
	Last24Hours Int64String `json:"1"`
}

type TickerFloatStats struct {
	Today       Float64String `json:"0"`
	Last24Hours Float64String `json:"1"`
}

type OHLC struct {
	ChannelID   Int64String `json:"0"`
	Data        OHLCData    `json:"1"`
	ChannelName string      `json:"2"`
	Pair        string      `json:"3"`
}

type TickerClose struct {
	Price     Float64String `json:"0"`
	LotVolume Float64String `json:"1"`
}

type OHLCData struct {
	Time                UnixTime      `json:"0"`
	EndTime             UnixTime      `json:"1"`
	Open                Float64String `json:"2"`
	High                Float64String `json:"3"`
	Low                 Float64String `json:"4"`
	Close               Float64String `json:"5"`
	VolumeWeightedPrice Float64String `json:"6"`
	Volume              Float64String `json:"7"`
	Count               int64         `json:"8"`
}

type Trade struct {
	ChannelID   int64       `json:"0"`
	Data        []TradeData `json:"1"`
	ChannelName string      `json:"2"`
	Pair        string      `json:"3"`
}

type TradeData struct {
	Price     Float64String `json:"0"`
	Volume    Float64String `json:"1"`
	Time      UnixTime      `json:"2"`
	Side      string        `json:"3"`
	OrderType string        `json:"4"`
	Misc      string        `json:"5"`
}

type Spread struct {
	ChannelID   int64      `json:"0"`
	Data        SpreadData `json:"1"`
	ChannelName string     `json:"2"`
	Pair        string     `json:"3"`
}

type SpreadData struct {
	Ask       Float64String `json:"0"`
	Bid       Float64String `json:"1"`
	Time      UnixTime      `json:"2"`
	BidVolume Float64String `json:"3"`
	AskVolume Float64String `json:"4"`
}

type Book struct {
	ChannelID   int64    `json:"0"`
	Data        BookData `json:"1"`
	ChannelName string   `json:"2"`
	Pair        string   `json:"3"`
}

type BookData struct {
	Asks []PriceLevel `json:"as"`
	Bids []PriceLevel `json:"bs"`
}

type BookUpdate struct {
	ChannelID   int64          `json:"0"`
	Data        BookUpdateData `json:"1"`
	ChannelName string         `json:"2"`
	Pair        string         `json:"3"`
}

type BookUpdateData struct {
	Asks []PriceLevel `json:"a"`
	Bids []PriceLevel `json:"b"`
}

type PriceLevel struct {
	Price     Float64String `json:"0"`
	Volume    Float64String `json:"1"`
	Timestamp UnixTime      `json:"2"`
}

type Error struct {
	Message string `json:"errorMessage"`
	Event   string `json:"event"`
	Status  string `json:"status"`
}
