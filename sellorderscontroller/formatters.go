package sellorderscontroller

import (
	"encoding/json"
	"time"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.SellOrder) ([]byte, error) {
	return json.MarshalIndent(toOrderResponse(order), "", "  ")
}

func formatGetResponse(order ordersservice.SellOrder) ([]byte, error) {
	return json.MarshalIndent(toOrderResponse(order), "", "  ")
}

func formatListResponse(orders []ordersservice.SellOrder) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = toOrderResponse(order)
	}

	return json.MarshalIndent(orderResponses, "", "  ")
}

func toOrderResponse(order ordersservice.SellOrder) _OrderResponse {
	return _OrderResponse{
		ID:        order.GetID(),
		Price:     float64(order.GetPrice()) / 1000,
		Quantity:  order.GetQuantity(),
		Ticker:    order.GetTicker(),
		Timestamp: order.GetTimestamp(),
	}
}

type _OrderResponse struct {
	ID        string    `json:"id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Ticker    string    `json:"ticker"`
	Timestamp time.Time `json:"timestamp"`
}
