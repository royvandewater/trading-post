package ordersservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// OrdersService manages CRUD for buy & sell orders
type OrdersService interface {
	CreateBuyOrder(buyOrder BuyOrder) (BuyOrder, int, error)
}

// New constructs a new OrdersService that will
// persist data using the provided mongo session
func New(mongoDB *mgo.Session) OrdersService {
	buyOrders := mongoDB.DB("tradingPost").C("buyOrders")
	return &service{buyOrders: buyOrders}
}

type service struct {
	buyOrders *mgo.Collection
}

func (s *service) CreateBuyOrder(buyOrder BuyOrder) (BuyOrder, int, error) {
	price, err := stockPrice(buyOrder.GetTicker())
	if err != nil {
		return nil, 502, err
	}

	buyOrder.SetPrice(price)
	err = s.buyOrders.Insert(buyOrder)
	if err != nil {
		return nil, 500, err
	}

	return buyOrder, 201, nil
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
