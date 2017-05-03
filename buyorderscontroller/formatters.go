package buyorderscontroller

import (
	"encoding/json"
	"time"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:        order.GetID(),
		Price:     order.GetPrice(),
		Quantity:  order.GetQuantity(),
		Ticker:    order.GetTicker(),
		Timestamp: order.GetTimestamp(),
	}, "", "  ")
}

func formatGetResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:        order.GetID(),
		Price:     order.GetPrice(),
		Quantity:  order.GetQuantity(),
		Ticker:    order.GetTicker(),
		Timestamp: order.GetTimestamp(),
	}, "", "  ")
}

func formatListResponse(orders []ordersservice.BuyOrder) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = _OrderResponse{
			ID:        order.GetID(),
			Price:     order.GetPrice(),
			Quantity:  order.GetQuantity(),
			Ticker:    order.GetTicker(),
			Timestamp: order.GetTimestamp(),
		}
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

type _OrderResponse struct {
	ID        string    `json:"id"`
	Price     float32   `json:"price"`
	Quantity  int       `json:"quantity"`
	Ticker    string    `json:"ticker"`
	Timestamp time.Time `json:"timestamp"`
}
