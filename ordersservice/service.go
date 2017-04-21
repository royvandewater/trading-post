package ordersservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/royvandewater/trading-post/usersservice"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrdersService manages CRUD for buy & sell orders
type OrdersService interface {
	// Create will persist a new order at the market rate
	// for the given ticker. The market rate will be subtracted
	// from the given user's riches.
	Create(userID, ticker string) (Order, error)

	// List returns all orders for a user
	List(userID string) ([]Order, error)
}

// New constructs a new OrdersService that will
// persist data using the provided mongo session
func New(mongoDB *mgo.Session, usersService usersservice.UsersService) OrdersService {
	orders := mongoDB.DB("tradingPost").C("orders")
	return &_Service{orders: orders, usersService: usersService}
}

type _Service struct {
	orders       *mgo.Collection
	usersService usersservice.UsersService
}

func (s *_Service) Create(userID, ticker string) (Order, error) {
	purchasePrice, err := stockPrice(ticker)
	if err != nil {
		return nil, err
	}

	order := NewOrder(userID, ticker, purchasePrice)
	err = s.orders.Insert(order)
	if err != nil {
		return nil, err
	}

	err = s.usersService.SubstractRichesByUserID(order.GetUserID(), purchasePrice)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *_Service) List(userID string) ([]Order, error) {
	var _orders []*_Order

	err := s.orders.Find(bson.M{"user_id": userID}).All(&_orders)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, len(_orders))
	for i, _order := range _orders {
		orders[i] = _order
	}

	return orders, nil
}

func stockPrice(ticker string) (float32, error) {
	url := fmt.Sprintf("https://stock.octoblu.com/last-trade-price/%v", ticker)
	response, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	if response.StatusCode != 200 {
		return 0, fmt.Errorf("Non 200 status code received from weather.octoblu.com: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	stockResponse := struct {
		Price float32 `json:"price"`
	}{}

	err = json.Unmarshal(data, &stockResponse)
	if err != nil {
		return 0, err
	}

	return stockResponse.Price, nil
}
