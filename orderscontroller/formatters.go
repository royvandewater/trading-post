package orderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.Order) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:            order.GetID(),
		PurchasePrice: order.GetPurchasePrice(),
		Ticker:        order.GetTicker(),
	}, "", "  ")
}

func formatListResponse(orders []ordersservice.Order) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = _OrderResponse{
			ID:            order.GetID(),
			PurchasePrice: order.GetPurchasePrice(),
			Ticker:        order.GetTicker(),
		}
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

type _OrderResponse struct {
	ID            string  `json:"id"`
	PurchasePrice float32 `json:"purchase_price"`
	Ticker        string  `json:"ticker"`
}
