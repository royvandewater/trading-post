package orderscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.Order) ([]byte, error) {
	return json.MarshalIndent(struct {
		PurchasePrice float32 `json:"purchasePrice"`
		Ticker        string  `json:"ticker"`
	}{
		Ticker:        order.GetTicker(),
		PurchasePrice: order.GetPurchasePrice(),
	}, "", "  ")
}
