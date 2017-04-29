package buyorderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:     order.GetID(),
		Price:  order.GetPrice(),
		Ticker: order.GetTicker(),
	}, "", "  ")
}

func formatGetResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:     order.GetID(),
		Price:  order.GetPrice(),
		Ticker: order.GetTicker(),
	}, "", "  ")
}

func formatListResponse(orders []ordersservice.BuyOrder) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = _OrderResponse{
			ID:     order.GetID(),
			Price:  order.GetPrice(),
			Ticker: order.GetTicker(),
		}
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

type _OrderResponse struct {
	ID     string  `json:"id"`
	Price  float32 `json:"price"`
	Ticker string  `json:"ticker"`
}
