package ordersservice

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// BuyOrder represents a purchase order for a
// specific stock
type BuyOrder interface {
	// GetID returns the unique identifier used to identify this
	// record in persistent storage
	GetID() string

	// GetPrice returns the price the order was purchased at
	// in increments of 1/1000 of a penny
	GetPrice() int

	// GetQuantity returns the quantity of stock purchased
	GetQuantity() int

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetTimestamp returns the timestamp from when the order was made
	GetTimestamp() time.Time

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewBuyOrder constructs a new order instance given a user, ticker, and price.
// it will also gain an ID, representing just this order instance
func NewBuyOrder(userID, ticker string, quantity, price int, timestamp time.Time) BuyOrder {
	return &_BuyOrder{
		ID:        uuid.NewV4().String(),
		Price:     price,
		Quantity:  quantity,
		Ticker:    ticker,
		Timestamp: timestamp.UTC(),
		UserID:    userID,
	}
}

type _BuyOrder struct {
	ID        string    `bson:"id"`
	Price     int       `bson:"price"`
	Quantity  int       `bson:"quantity"`
	Ticker    string    `bson:"ticker"`
	Timestamp time.Time `bson:"timestamp"`
	UserID    string    `bson:"user_id"`
}

func (order *_BuyOrder) GetID() string {
	return order.ID
}

func (order *_BuyOrder) GetPrice() int {
	return order.Price
}

func (order *_BuyOrder) GetQuantity() int {
	return order.Quantity
}

func (order *_BuyOrder) GetTicker() string {
	return order.Ticker
}

func (order *_BuyOrder) GetTimestamp() time.Time {
	return order.Timestamp
}

func (order *_BuyOrder) GetUserID() string {
	return order.UserID
}
