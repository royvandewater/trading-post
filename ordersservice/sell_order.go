package ordersservice

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// SellOrder represents a sell order for a
// specific stock
type SellOrder interface {
	// GetID returns the unique identifier used to identify this
	// record in persistent storage
	GetID() string

	// GetPrice returns the price the order was purchased at
	// in increments of 1/1000 of a penny
	GetPrice() int

	// GetQuantity returns the quantity of stock sold
	GetQuantity() int

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetTimestamp returns the timestamp from when the order was created
	GetTimestamp() time.Time

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewSellOrder constructs a new order instance given a user, ticker, and price.
// it will also gain an ID, representing just this order instance
func NewSellOrder(userID, ticker string, quantity, price int, timestamp time.Time) SellOrder {
	return &_SellOrder{
		ID:        uuid.NewV4().String(),
		Price:     price,
		Quantity:  quantity,
		Ticker:    ticker,
		Timestamp: timestamp.UTC(),
		UserID:    userID,
	}
}

type _SellOrder struct {
	ID        string    `bson:"id"`
	Price     int       `bson:"price"`
	Quantity  int       `bson:"quantity"`
	Ticker    string    `bson:"ticker"`
	Timestamp time.Time `bson:"timestamp"`
	UserID    string    `bson:"user_id"`
}

func (order *_SellOrder) GetID() string {
	return order.ID
}

func (order *_SellOrder) GetPrice() int {
	return order.Price
}

func (order *_SellOrder) GetQuantity() int {
	return order.Quantity
}

func (order *_SellOrder) GetTicker() string {
	return order.Ticker
}

func (order *_SellOrder) GetTimestamp() time.Time {
	return order.Timestamp
}

func (order *_SellOrder) GetUserID() string {
	return order.UserID
}
