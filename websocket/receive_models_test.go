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
			expectedModel: &Pong{Event: "pong", ReqID: 42},
			expectedError: nil,
		},
		{
			name:          "heartbeat",
			bytes:         []byte(`{"event": "heartbeat"}`),
			expectedModel: &HeartBeat{Event: "heartbeat"},
			expectedError: nil,
		},
		{
			name:  "systemStatus",
			bytes: []byte(`{"connectionID":17978356104855020991,"event":"systemStatus","status":"online","version":"1.5.1"}`),
			expectedModel: &SystemStatus{
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
			expectedModel: &SubscriptionStatus{
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
			expectedModel: &Ticker{
				ChannelID: 916,
				Data: TickerData{
					Ask:                   TickerAskBid{Price: 0.427, WholeLotVolume: 16169, LotVolume: 16169.08316400},
					Bid:                   TickerAskBid{Price: 0.4269, WholeLotVolume: 1000, LotVolume: 1000.0},
					Close:                 TickerClose{Price: 0.427, LotVolume: 270.856836},
					Volume:                TickerFloatStats{Today: 57719824.25617952, Last24Hours: 60354910.83998816},
					VolumeWeightedAverage: TickerFloatStats{Today: 0.40286360, Last24Hours: 0.40226657},
					Trades:                TickerTrades{Today: 22509, Last24Hours: 23886},
					Low:                   TickerFloatStats{Today: 0.36, Last24Hours: 0.36},
					High:                  TickerFloatStats{Today: 0.43605, Last24Hours: 0.43605},
					Open:                  TickerFloatStats{Today: 0.38529, Last24Hours: 0.39564},
				},
				ChannelName: "ticker",
				Pair:        "XBT/EUR",
			},
			expectedError: nil,
		},
		{
			name:  "ohlc",
			bytes: []byte(`[920,["1608207638.842716","1608207900.000000","0.46228000","0.46256000","0.46200000","0.46256000","0.46210203","86513.88714914",15],"ohlc-5","XBT/EUR"]`),
			expectedModel: &OHLC{
				ChannelID: 920,
				Data: OHLCData{
					Time:                time.Unix(1608207638, 842715978),
					EndTime:             time.Unix(1608207900, 0),
					Open:                0.46228,
					High:                0.46256,
					Low:                 0.462,
					Close:               0.46256,
					VolumeWeightedPrice: 0.46210203,
					Volume:              86513.88714914,
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
			expectedModel: &Trade{
				ChannelID:   0,
				ChannelName: "trade",
				Pair:        "XBT/USD",
				Data: []TradeData{
					{
						Price:     5541.2,
						Volume:    0.15850568,
						Time:      time.Unix(1534614057, 321597099),
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
			expectedModel: &Spread{
				ChannelID:   0,
				ChannelName: "spread",
				Pair:        "XBT/USD",
				Data: SpreadData{
					Ask:       5698.4,
					Bid:       5700,
					Time:      time.Unix(1542057299, 545897006),
					BidVolume: 1.01234567,
					AskVolume: 0.98765432,
				},
			},
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
