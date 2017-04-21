package ordersservice

import uuid "github.com/satori/go.uuid"

// Order represents a purchase order for a
// specific stock
type Order interface {
	// GetID returns the unique identifier used to identify this
	// record in persistent storage
	GetID() string

	// GetPurchasePrice returns the price the order was purchased at
	GetPurchasePrice() float32

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewOrder constructs a new order instance given a user, ticker, and price.
// it will also gain an ID, representing just this order instance
func NewOrder(userID, ticker string, purchasePrice float32) Order {
	return &_Order{
		ID:            uuid.NewV4().String(),
		PurchasePrice: purchasePrice,
		Ticker:        ticker,
		UserID:        userID,
	}
}

type _Order struct {
	PurchasePrice float32 `bson:"purchase_price"`
	Ticker        string  `bson:"ticker"`
	UserID        string  `bson:"user_id"`
	ID            string  `bson:"id"`
}

func (order *_Order) GetID() string {
	return order.ID
}

func (order *_Order) GetPurchasePrice() float32 {
	return order.PurchasePrice
}

func (order *_Order) GetTicker() string {
	return order.Ticker
}

func (order *_Order) GetUserID() string {
	return order.UserID
}
