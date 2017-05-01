package ordersservice

import uuid "github.com/satori/go.uuid"

// SellOrder represents a sell order for a
// specific stock
type SellOrder interface {
	// GetID returns the unique identifier used to identify this
	// record in persistent storage
	GetID() string

	// GetPrice returns the price the order was purchased at
	GetPrice() float32

	// GetQuantity returns the quantity of stock sold
	GetQuantity() int

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewSellOrder constructs a new order instance given a user, ticker, and price.
// it will also gain an ID, representing just this order instance
func NewSellOrder(userID, ticker string, quantity int, price float32) SellOrder {
	return &_SellOrder{
		ID:       uuid.NewV4().String(),
		Price:    price,
		Quantity: quantity,
		Ticker:   ticker,
		UserID:   userID,
	}
}

type _SellOrder struct {
	ID       string  `bson:"id"`
	Price    float32 `bson:"price"`
	Quantity int     `bson:"quantity"`
	Ticker   string  `bson:"ticker"`
	UserID   string  `bson:"user_id"`
}

func (order *_SellOrder) GetID() string {
	return order.ID
}

func (order *_SellOrder) GetPrice() float32 {
	return order.Price
}

func (order *_SellOrder) GetQuantity() int {
	return order.Quantity
}

func (order *_SellOrder) GetTicker() string {
	return order.Ticker
}

func (order *_SellOrder) GetUserID() string {
	return order.UserID
}
