package websocket

import (
	"encoding/json"
	"time"
)

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

type TickerTrades struct {
	Today       int64
	Last24Hours int64
}

type TickerFloatStats struct {
	Today       float64
	Last24Hours float64
}

type OHLC struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        OHLCData
}

type TickerClose struct {
	Price     float64
	LotVolume float64
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

type Book struct {
	ChannelID   int64
	ChannelName string
	Pair        string
	Data        BookData
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

type Error struct {
	Message string `json:"errorMessage"`
	Event   string `json:"event"`
	Status  string `json:"status"`
}
