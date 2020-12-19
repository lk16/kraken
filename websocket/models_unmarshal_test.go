package websocket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalReceivedMessage(t *testing.T) {

	type testCase struct {
		name          string
		bytes         []byte
		expectedModel interface{}
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "pong",
			bytes:         []byte(`{"event": "pong", "ReqID": 42}`),
			expectedModel: Pong{Event: "pong", ReqID: 42},
			expectedError: nil,
		},
		{
			name:          "heartbeat",
			bytes:         []byte(`{"event": "heartbeat"}`),
			expectedModel: HeartBeat{Event: "heartbeat"},
			expectedError: nil,
		},
		{
			name:  "systemStatus",
			bytes: []byte(`{"connectionID":17978356104855020991,"event":"systemStatus","status":"online","version":"1.5.1"}`),
			expectedModel: SystemStatus{
				ConnectionID: "17978356104855020991",
				Event:        "systemStatus",
				Status:       "online",
				Version:      "1.5.1",
			},
			expectedError: nil,
		},
		{
			name:  "subscriptionStatus",
			bytes: []byte(`{"channelID":916,"channelName":"ticker","event":"subscriptionStatus","pair":"XBT/EUR","status":"subscribed","subscription":{"name":"ticker"}}`),
			expectedModel: SubscriptionStatus{
				ChannelID:   916,
				ChannelName: "ticker",
				Event:       "subscriptionStatus",
				Pair:        "XBT/EUR",
				Status:      "subscribed",
				Subscription: SubscriptionDetails{
					Name: "ticker",
				},
			},
			expectedError: nil,
		},
		{
			name: "ticker",
			bytes: []byte(`[916,{"a":["0.42700000",16169,"16169.08316400"],"b":["0.42690000",1000,"1000.00000000"],"c":` +
				`["0.42700000","270.85683600"],"v":["57719824.25617952","60354910.83998816"],"p":["0.40286360","0.40226657"]` +
				`,"t":[22509,23886],"l":["0.36000000","0.36000000"],"h":["0.43605000","0.43605000"],"o":["0.38529000","0.395` +
				`64000"]},"ticker","XBT/EUR"]`),
			expectedModel: Ticker{
				ChannelID: Int64String(916),
				Data: TickerData{
					Ask: TickerAskBid{
						Price:          Float64String(0.427),
						WholeLotVolume: Int64String(16169),
						LotVolume:      Float64String(16169.08316400),
					},
					Bid: TickerAskBid{
						Price:          Float64String(0.4269),
						WholeLotVolume: Int64String(1000),
						LotVolume:      Float64String(1000.0),
					},
					Close: TickerClose{Price: 0.427, LotVolume: 270.856836},
					Volume: TickerFloatStats{
						Today:       Float64String(57719824.25617952),
						Last24Hours: Float64String(60354910.83998816),
					},

					VolumeWeightedAverage: TickerFloatStats{
						Today:       Float64String(0.40286360),
						Last24Hours: Float64String(0.40226657),
					},
					Trades: TickerTrades{Today: Int64String(22509), Last24Hours: Int64String(23886)},
					Low: TickerFloatStats{
						Today:       Float64String(0.36),
						Last24Hours: Float64String(0.36),
					},
					High: TickerFloatStats{
						Today:       Float64String(0.43605),
						Last24Hours: Float64String(0.43605),
					},
					Open: TickerFloatStats{
						Today:       Float64String(0.38529),
						Last24Hours: Float64String(0.39564),
					},
				},
				ChannelName: "ticker",
				Pair:        "XBT/EUR",
			},
			expectedError: nil,
		},
		{
			name:  "ohlc",
			bytes: []byte(`[920,["1608207638.842716","1608207900.000000","0.46228000","0.46256000","0.46200000","0.46256000","0.46210203","86513.88714914",15],"ohlc-5","XBT/EUR"]`),
			expectedModel: OHLC{
				ChannelID: 920,
				Data: OHLCData{
					Time:                UnixTime(time.Unix(1608207638, 842715978)),
					EndTime:             UnixTime(time.Unix(1608207900, 0)),
					Open:                Float64String(0.46228),
					High:                Float64String(0.46256),
					Low:                 Float64String(0.462),
					Close:               Float64String(0.46256),
					VolumeWeightedPrice: Float64String(0.46210203),
					Volume:              Float64String(86513.88714914),
					Count:               15,
				},
				ChannelName: "ohlc-5",
				Pair:        "XBT/EUR",
			},
			expectedError: nil,
		},
		{
			name:  "trade",
			bytes: []byte(`[0,[["5541.20000","0.15850568","1534614057.321597","s","l","foo"]],"trade","XBT/USD"]`),
			expectedModel: Trade{
				ChannelID:   0,
				ChannelName: "trade",
				Pair:        "XBT/USD",
				Data: []TradeData{
					{
						Price:     5541.2,
						Volume:    0.15850568,
						Time:      UnixTime(time.Unix(1534614057, 321597099)),
						Side:      "s",
						OrderType: "l",
						Misc:      "foo",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:  "spread",
			bytes: []byte(`[0,["5698.40000","5700.00000","1542057299.545897","1.01234567","0.98765432"],"spread","XBT/USD"]`),
			expectedModel: Spread{
				ChannelID:   0,
				ChannelName: "spread",
				Pair:        "XBT/USD",
				Data: SpreadData{
					Ask:       5698.4,
					Bid:       5700,
					Time:      UnixTime(time.Unix(1542057299, 545897006)),
					BidVolume: 1.01234567,
					AskVolume: 0.98765432,
				},
			},
			expectedError: nil,
		},
		{
			name:  "bookSnapshot",
			bytes: []byte(`[0,{"as":[["5541.30000","2.50700000","1534614248.123678"]],"bs":[["5541.20000","1.52900000","1534614248.765567"]]},"book-100","XBT/USD"]`),
			expectedModel: Book{
				ChannelID:   0,
				ChannelName: "book-100",
				Pair:        "XBT/USD",
				Data: BookData{
					Asks: []PriceLevel{
						{
							Price:     5541.3,
							Volume:    2.507,
							Timestamp: UnixTime(time.Unix(1534614248, 123677968)),
						},
					},
					Bids: []PriceLevel{
						{
							Price:     5541.2,
							Volume:    1.529,
							Timestamp: UnixTime(time.Unix(1534614248, 765567064)),
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:  "bookUpdate",
			bytes: []byte(`[1234,{"a":[["5541.30000","2.50700000","1534614248.456738"]],"c":"974942666"},"book-10","XBT/USD"]`),
			expectedModel: BookUpdate{
				ChannelID:   1234,
				ChannelName: "book-10",
				Pair:        "XBT/USD",
				Data: BookUpdateData{
					Asks: []PriceLevel{
						{
							Price:     5541.3,
							Volume:    2.507,
							Timestamp: UnixTime(time.Unix(1534614248, 456737995)),
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:  "bookUpdateWithAsksAndBids",
			bytes: []byte(`[912,{"a":[["0.46800000","2940.56589429","1608240638.875519"]]},{"b":[["0.46877000","0.00000000","1608240638.875818"]],"c":"751501448"},"book-10","XRP/EUR"]`),
			expectedModel: BookUpdate{
				ChannelID:   912,
				ChannelName: "book-10",
				Pair:        "XRP/EUR",
				Data: BookUpdateData{
					Asks: []PriceLevel{
						{
							Price:     0.468,
							Volume:    2940.56589429,
							Timestamp: UnixTime(time.Unix(1608240638, 875519037)),
						},
					},
					Bids: []PriceLevel{
						{
							Price:     0.46877,
							Volume:    0.0,
							Timestamp: UnixTime(time.Unix(1608240638, 875818014)),
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:  "error",
			bytes: []byte(`{"errorMessage": "EOrder:Order minimum not met", "event": "addOrderStatus", "status": "error"}`),
			expectedModel: Error{
				Message: "EOrder:Order minimum not met",
				Event:   "addOrderStatus",
				Status:  "error",
			},
			expectedError: nil,
		},
		{
			name: "ownTrades",
			bytes: []byte(`[[{"3HA3PV-3HA3P-3HA3PV":{"cost":"99.90000737","fee":"0.09792001","margin":"0.00000000","ordertxid":"ORWUAU-YUFW6-PSGDHG","ordertype":"limit","pair":"XRP/EUR",
				"postxid":"TKH2SE-M7IF5-CFI7LT","price":"0.47957000","time":"1237535943.237535","type":"sell","vol":"123.456789"}}],"ownTrades",{"sequence":1}]`),
			expectedModel: nil,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model, err := unmarshalReceivedMessage(testCase.bytes)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedModel, model)
		})
	}
}
