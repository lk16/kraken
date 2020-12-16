package websocket

import (
	"testing"

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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model, err := unmarshalReceivedMessage(testCase.bytes)

			assert.Equal(t, testCase.expectedModel, model)
			assert.Equal(t, testCase.expectedError, err)
		})
	}

}
