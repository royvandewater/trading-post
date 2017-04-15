package ordersservice

import (
	mgo "gopkg.in/mgo.v2"
)

// OrdersService manages CRUD for buy & sell orders
type OrdersService interface {
	CreateBuyOrder(buyOrder BuyOrder) (BuyOrder, int, error)
}

// New constructs a new OrdersService instance using the
// provided firebase auth token
func New(mongoDB *mgo.Session) OrdersService {
	buyOrders := mongoDB.DB("tradingPost").C("buyOrders")
	return &service{buyOrders: buyOrders}
}

type service struct {
	buyOrders *mgo.Collection
}

func (s *service) CreateBuyOrder(buyOrder BuyOrder) (BuyOrder, int, error) {
	err := s.buyOrders.Insert(buyOrder)
	if err != nil {
		return nil, 500, err
	}

	return buyOrder, 201, nil
}
