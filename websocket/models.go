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
	ChannelID   Int64String
	Data        TickerData
	ChannelName string
	Pair        string
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
	Price          Float64String
	WholeLotVolume Int64String
	LotVolume      Float64String
}

type TickerTrades struct {
	Today       Int64String
	Last24Hours Int64String
}

type TickerFloatStats struct {
	Today       Float64String
	Last24Hours Float64String
}

type OHLC struct {
	ChannelID   Int64String
	Data        OHLCData
	ChannelName string
	Pair        string
}

type TickerClose struct {
	Price     Float64String
	LotVolume Float64String
}

type OHLCData struct {
	Time                UnixTime
	EndTime             UnixTime
	Open                Float64String
	High                Float64String
	Low                 Float64String
	Close               Float64String
	VolumeWeightedPrice Float64String
	Volume              Float64String
	Count               int64
}

type Trade struct {
	ChannelID   int64
	Data        []TradeData
	ChannelName string
	Pair        string
}

type TradeData struct {
	Price     Float64String
	Volume    Float64String
	Time      UnixTime
	Side      string
	OrderType string
	Misc      string
}

type Spread struct {
	ChannelID   int64
	Data        SpreadData
	ChannelName string
	Pair        string
}

type SpreadData struct {
	Ask       Float64String
	Bid       Float64String
	Time      UnixTime
	BidVolume Float64String
	AskVolume Float64String
}

type Book struct {
	ChannelID   int64
	Data        BookData
	ChannelName string
	Pair        string
}

type BookData struct {
	Asks []PriceLevel `json:"as"`
	Bids []PriceLevel `json:"bs"`
}

type BookUpdate struct {
	ChannelID   int64
	Data        BookUpdateData
	ChannelName string
	Pair        string
}

type BookUpdateData struct {
	Asks []PriceLevel `json:"a"`
	Bids []PriceLevel `json:"b"`
}

type PriceLevel struct {
	Price     Float64String
	Volume    Float64String
	Timestamp UnixTime
}

type Error struct {
	Message string `json:"errorMessage"`
	Event   string `json:"event"`
	Status  string `json:"status"`
}
