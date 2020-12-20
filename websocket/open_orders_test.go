package websocket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testOrder = OpenOrders{
	Orders: map[string]OpenOrder{
		"OGTT3Y-C6I3P-XRI6HX": {
			Cost: 1.00000,
			Description: OpenOrderDescription{
				ConditionalClose: "",
				Leverage:         "0:1",
				Order:            "sell 10.00345345 XBT/EUR @ limit 34.50000 with 0:1 leverage",
				OrderType:        "limit",
				Pair:             "XBT/EUR",
				Price:            34.50000,
				Price2:           55.00000,
				Type:             "sell",
			},
			ExpirationTime: UnixTime(time.Unix(0, 0)),
			Fee:            0.00000,
			LimitPrice:     34.50000,
			Miscellaneous:  "",
			OFlags:         "fcib",
			OpenTime:       UnixTime(time.Unix(0, 0)),
			Price:          34.50000,
			ReferenceID:    "OKIVMP-5GVZN-Z2D2UA",
			StartTime:      UnixTime(time.Unix(0, 0)),
			Status:         "pending",
			StopPrice:      0.000000,
			UserReference:  0,
			Volume:         10.00345345,
			VolumeExecuted: 9.00000000,
		},
	},
}

func TestOpenOrdersUpdate(t *testing.T) {

	t.Run("initialUpdate", func(t *testing.T) {
		var state OpenOrders
		state.Update(testOrder)
		assert.Equal(t, testOrder, state)
	})

	t.Run("updateOrderStatus", func(t *testing.T) {
		state := testOrder

		update := OpenOrders{
			Orders: map[string]OpenOrder{
				"OGTT3Y-C6I3P-XRI6HX": {Status: "open"},
			},
		}

		state.Update(update)

		// updating field of map value in golang can't be done directly
		expected := testOrder
		order := expected.Orders["OGTT3Y-C6I3P-XRI6HX"]
		order.Status = "open"
		expected.Orders["OGTT3Y-C6I3P-XRI6HX"] = order

		assert.Equal(t, expected, state)
	})

	t.Run("fillUpdate", func(t *testing.T) {
		state := testOrder

		//[[{"OLFCT6-43DXW-EABUVM":{"cost":"30.00000163","vol_exec":"70.90020000","fee":"0.04800000","avg_price":"0.42313000"}}],"openOrders",{"sequence":4}]

		update := OpenOrders{
			Orders: map[string]OpenOrder{
				"OGTT3Y-C6I3P-XRI6HX": {
					Cost:           30.00000163,
					VolumeExecuted: 70.90020000,
					Fee:            0.04800000,
					AveragePrice:   0.42313000,
				},
			},
		}

		state.Update(update)

		// updating field of map value in golang can't be done directly
		expected := testOrder
		order := expected.Orders["OGTT3Y-C6I3P-XRI6HX"]
		order.Cost = 30.00000163
		order.VolumeExecuted = 70.90020000
		order.Fee = 0.04800000
		order.AveragePrice = 0.42313000
		expected.Orders["OGTT3Y-C6I3P-XRI6HX"] = order

		assert.Equal(t, expected, state)

		update = OpenOrders{
			Orders: map[string]OpenOrder{
				"OGTT3Y-C6I3P-XRI6HX": {
					Status:         "closed",
					Cost:           30.00000163,
					VolumeExecuted: 70.90020000,
					Fee:            0.04800000,
					AveragePrice:   0.42313000,
				},
			},
		}

		state.Update(update)

		order = expected.Orders["OGTT3Y-C6I3P-XRI6HX"]
		order.Status = "closed"
		expected.Orders["OGTT3Y-C6I3P-XRI6HX"] = order

		assert.Equal(t, expected, state)
	})

	t.Run("cancelUpdate", func(t *testing.T) {
		state := testOrder

		update := OpenOrders{
			Orders: map[string]OpenOrder{
				"OGTT3Y-C6I3P-XRI6HX": {
					Status:         "canceled",
					Cost:           0.00000000,
					VolumeExecuted: 0.00000000,
					Fee:            0.00000000,
					AveragePrice:   0.00000000,
				},
			},
		}

		state.Update(update)

		// updating field of map value in golang can't be done directly
		expected := testOrder
		order := expected.Orders["OGTT3Y-C6I3P-XRI6HX"]
		order.Status = "closed"
		expected.Orders["OGTT3Y-C6I3P-XRI6HX"] = order

		assert.Equal(t, expected, state)
	})
}
