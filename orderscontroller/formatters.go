package orderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.Order) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		Ticker:        order.GetTicker(),
		PurchasePrice: order.GetPurchasePrice(),
	}, "", "  ")
}

func formatListResponse(orders []ordersservice.Order) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = _OrderResponse{
			Ticker:        order.GetTicker(),
			PurchasePrice: order.GetPurchasePrice(),
		}
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

type _OrderResponse struct {
	PurchasePrice float32 `json:"purchase_price"`
	Ticker        string  `json:"ticker"`
}
