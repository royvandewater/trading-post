package sellorderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.SellOrder) ([]byte, error) {
	return json.MarshalIndent(_OrderResponse{
		ID:     order.GetID(),
		Price:  order.GetPrice(),
		Ticker: order.GetTicker(),
	}, "", "  ")
}

type _OrderResponse struct {
	ID     string  `json:"id"`
	Price  float32 `json:"price"`
	Ticker string  `json:"ticker"`
}
