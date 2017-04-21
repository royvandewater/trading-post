package ordersservice

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Order represents a purchase order for a
// specific stock
type Order interface {
	// JSON serializes the buy order
	JSON() ([]byte, error)

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewOrder constructs a new order instance given a user, ticker, and price
func NewOrder(userID, ticker string, purchasePrice float32) Order {
	return &_Order{
		UserID:        userID,
		Ticker:        ticker,
		PurchasePrice: purchasePrice,
	}
}

// ParseOrderForUserID instantiates a new Order instance from JSON data
func ParseOrderForUserID(userID string, data io.ReadCloser) (Order, error) {
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	order := &_Order{}
	err = json.Unmarshal(dataBytes, order)
	if err != nil {
		return nil, err
	}

	order.UserID = userID

	return order, nil
}

type _Order struct {
	PurchasePrice float32 `json:"purchase_price"`
	Ticker        string  `json:"ticker"`
	UserID        string  `json:"user_id"`
}

func (order *_Order) JSON() ([]byte, error) {
	return json.MarshalIndent(order, "", "  ")
}

func (order *_Order) GetTicker() string {
	return order.Ticker
}

func (order *_Order) GetUserID() string {
	return order.UserID
}
