package ordersservice

// Order represents a purchase order for a
// specific stock
type Order interface {
	// GetPurchasePrice returns the price the order was purchased at
	GetPurchasePrice() float32

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

type _Order struct {
	PurchasePrice float32 `bson:"purchase_price"`
	Ticker        string  `bson:"ticker"`
	UserID        string  `bson:"user_id"`
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
