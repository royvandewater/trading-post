package sellorderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.SellOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:       order.GetID(),
		Price:    order.GetPrice(),
		Quantity: order.GetQuantity(),
		Ticker:   order.GetTicker(),
	}, "", "  ")
}

func formatGetResponse(order ordersservice.SellOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:       order.GetID(),
		Price:    order.GetPrice(),
		Quantity: order.GetQuantity(),
		Ticker:   order.GetTicker(),
	}, "", "  ")
}

func formatListResponse(orders []ordersservice.SellOrder) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = _OrderResponse{
			ID:       order.GetID(),
			Price:    order.GetPrice(),
			Quantity: order.GetQuantity(),
			Ticker:   order.GetTicker(),
		}
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

type _OrderResponse struct {
	ID       string  `json:"id"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Ticker   string  `json:"ticker"`
}
