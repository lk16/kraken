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

	if string(bytes) == "null" {
		return nil
	}

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

	if len(bytes) == 0 {
		return nil
	}

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
	Pair         []string     `json:"pair,omitempty"`
	Subscription Subscription `json:"subscription"`
}

type Unsubscribe struct {
	Event        string       `json:"event"`
	ReqID        int          `json:"reqid,omitempty"`
	Pair         []string     `json:"pair,omitempty"`
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

type OwnTrades struct {
	Trades      []map[string]OwnTrade
	ChannelName string
	Sequence    Sequence
}

type OwnTrade struct {
	Cost               Float64String `json:"cost"`
	Fee                Float64String `json:"fee"`
	Margin             Float64String `json:"margin"`
	OrderTransactionID string        `json:"ordertxid"`
	OrderType          string        `json:"ordertype"`
	Pair               string        `json:"pair"`
	PosTransactionID   string        `json:"postxid"`
	Price              Float64String `json:"price"`
	Time               UnixTime      `json:"time"`
	Type               string        `json:"type"`
	Volume             Float64String `json:"vol"`
}

type Sequence struct {
	Sequence int64 `json:"sequence"`
}

type OpenOrders struct {
	Orders      map[string]OpenOrder
	ChannelName string
	Sequence    Sequence
}

type OpenOrder struct {
	Cost           Float64String        `json:"cost"`
	Description    OpenOrderDescription `json:"descr"`
	ExpirationTime UnixTime             `json:"expiretm"`
	Fee            Float64String        `json:"fee"`
	LimitPrice     Float64String        `json:"limitprice"`
	Miscellaneous  string               `json:"misc"`
	OFlags         string               `json:"oflags"`
	OpenTime       UnixTime             `json:"opentm"`
	Price          Float64String        `json:"price"`
	ReferenceID    string               `json:"refid"`
	StartTime      UnixTime             `json:"starttm"`
	Status         string               `json:"status"`
	StopPrice      Float64String        `json:"stopprice"`
	UserReference  int64                `json:"userref"`
	Volume         Float64String        `json:"vol"`
	VolumeExecuted Float64String        `json:"vol_exec"`
	AveragePrice   Float64String        `json:"avg_price"`
	CancelReason   string               `json:"cancel_reason"`
}

type OpenOrderDescription struct {
	ConditionalClose string        `json:"close"`
	Leverage         string        `json:"leverage"`
	Order            string        `json:"order"`
	OrderType        string        `json:"ordertype"`
	Pair             string        `json:"pair"`
	Price            Float64String `json:"price"`
	Price2           Float64String `json:"price2"`
	Type             string        `json:"type"`
}

type AddOrder struct {
	Event            string `json:"event"`
	Token            string `json:"token"`
	ReqID            int64  `json:"reqid"`
	OrderType        string `json:"ordertype"`
	Type             string `json:"type"`
	Pair             string `json:"pair"`
	Price            string `json:"price"`
	Price2           string `json:"price2"`
	Volume           string `json:"volume"`
	Leverage         string `json:"leverage"`
	OFlags           string `json:"oflags"`
	StartTime        string `json:"starttm"`
	ExpireTime       string `json:"expiretm"`
	UserReference    string `json:"userref"`
	Validate         string `json:"validate"`
	CloseOrderType   string `json:"close[ordertype]"`
	ClosePrice       string `json:"close[price]"`
	ClosePrice2      string `json:"close[price2]"`
	TradingAgreement string `json:"trading_agreement"`
}

type CancelOrder struct {
	Event         string   `json:"event"`
	Token         string   `json:"token"`
	ReqID         int      `json:"reqid"`
	TransactionID []string `json:"txid"`
}

type CancelOrderStatus struct {
	Event        string `json:"event"`
	ReqID        int    `json:"reqid"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type CancelAll struct {
	Event string `json:"event"`
	Token string `json:"token"`
}
