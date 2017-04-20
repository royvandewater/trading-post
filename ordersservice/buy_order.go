package ordersservice

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// BuyOrder represents a purchase order for a
// specific stock
type BuyOrder interface {
	// JSON serializes the buy order
	JSON() ([]byte, error)

	// SetPrice stores the price of the stock
	SetPrice(price float32)

	// GetTicker returns the stock ticker code
	GetTicker() string
}

// ParseBuyOrderForUserID instantiates a new BuyOrder instance from JSON
// data
func ParseBuyOrderForUserID(userID string, data io.ReadCloser) (BuyOrder, error) {
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	buyOrder := &buyOrder{}
	err = json.Unmarshal(dataBytes, buyOrder)
	if err != nil {
		return nil, err
	}

	buyOrder.UserID = userID

	return buyOrder, nil
}

type buyOrder struct {
	Price  float32 `json:"price"`
	Ticker string  `json:"ticker"`
	UserID string  `json:"user_id"`
}

func (order *buyOrder) JSON() ([]byte, error) {
	return json.Marshal(order)
}

func (order *buyOrder) SetPrice(price float32) {
	order.Price = price
}

func (order *buyOrder) GetTicker() string {
	return order.Ticker
}
