package websocket

func (openOrders *OpenOrders) DeleteInactiveOrders() {
	for orderID, order := range openOrders.Orders {
		if order.Status == "canceled" || order.Status == "closed" {
			delete(openOrders.Orders, orderID)
		}
	}
}

func (openOrders *OpenOrders) Update(update OpenOrders) {
	for orderID, orderUpdate := range update.Orders {
		openOrders.updateOrder(orderID, orderUpdate)
	}
}

func (openOrders *OpenOrders) updateOrder(orderID string, update OpenOrder) {

	current, ok := openOrders.Orders[orderID]
	if !ok {
		if openOrders.Orders == nil {
			openOrders.Orders = make(map[string]OpenOrder)
		}

		openOrders.Orders[orderID] = update
		return
	}

	var zeroValue OpenOrder

	if update.Cost != zeroValue.Cost {
		current.Cost = update.Cost
	}

	if update.Description != zeroValue.Description {
		current.Description = update.Description
	}

	if update.ExpirationTime != zeroValue.ExpirationTime {
		current.ExpirationTime = update.ExpirationTime
	}

	if update.Fee != zeroValue.Fee {
		current.Fee = update.Fee
	}

	if update.LimitPrice != zeroValue.LimitPrice {
		current.LimitPrice = update.LimitPrice
	}

	if update.Miscellaneous != zeroValue.Miscellaneous {
		current.Miscellaneous = update.Miscellaneous
	}

	if update.OFlags != zeroValue.OFlags {
		current.OFlags = update.OFlags
	}

	if update.OpenTime != zeroValue.OpenTime {
		current.OpenTime = update.OpenTime
	}

	if update.Price != zeroValue.Price {
		current.Price = update.Price
	}

	if update.ReferenceID != zeroValue.ReferenceID {
		current.ReferenceID = update.ReferenceID
	}

	if update.StartTime != zeroValue.StartTime {
		current.StartTime = update.StartTime
	}

	if update.Status != zeroValue.Status {
		current.Status = update.Status
	}

	if update.StopPrice != zeroValue.StopPrice {
		current.StopPrice = update.StopPrice
	}

	if update.UserReference != zeroValue.UserReference {
		current.UserReference = update.UserReference
	}

	if update.Volume != zeroValue.Volume {
		current.Volume = update.Volume
	}

	if update.VolumeExecuted != zeroValue.VolumeExecuted {
		current.VolumeExecuted = update.VolumeExecuted
	}

	if update.AveragePrice != zeroValue.AveragePrice {
		current.AveragePrice = update.AveragePrice
	}

	if update.CancelReason != zeroValue.CancelReason {
		current.CancelReason = update.CancelReason

	}

	openOrders.Orders[orderID] = current
}
