package ordersservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/royvandewater/trading-post/usersservice"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrdersService manages CRUD for buy & sell orders
type OrdersService interface {
	// CreateBuyOrder will persist a new order at the market rate
	// for the given ticker. The market rate will be subtracted
	// from the given profile's riches and the user's stock
	// quantity for this ticker will go up.
	CreateBuyOrder(userID, ticker string, quantity int) (BuyOrder, error)

	// GetBuyOrder will retrieve a buy order for a given user & id
	GetBuyOrder(userID, id string) (BuyOrder, error)

	// GetSellOrder will retrieve a sell order for a given user & id
	GetSellOrder(userID, id string) (SellOrder, error)

	// ListBuyOrders returns all orders for a user
	ListBuyOrders(userID string) ([]BuyOrder, error)

	// ListSellOrders returns all orders for a user
	ListSellOrders(userID string) ([]SellOrder, error)

	// CreateSellOrder will persist a new sell order at the market rate
	// for the given ticker. The market rate will be added
	// to the given user's riches and the user's stock quantity
	// for this ticker will go down. It will return an error if
	// the user doesn't own enough of the stock
	CreateSellOrder(userID, ticker string, quantity int) (SellOrder, error)
}

// New constructs a new BuyOrdersService that will
// persist data using the provided mongo session
func New(mongoDB *mgo.Session, usersService usersservice.UsersService) OrdersService {
	buyOrders := mongoDB.DB("trading_post").C("buy_orders")
	sellOrders := mongoDB.DB("trading_post").C("sell_orders")
	return &_Service{
		buyOrders:    buyOrders,
		sellOrders:   sellOrders,
		usersService: usersService,
	}
}

type _Service struct {
	buyOrders    *mgo.Collection
	sellOrders   *mgo.Collection
	usersService usersservice.UsersService
}

func (s *_Service) CreateBuyOrder(userID, ticker string, quantity int) (BuyOrder, error) {
	purchasePrice, err := stockPrice(ticker)
	if err != nil {
		return nil, err
	}
	if purchasePrice <= 0 {
		return nil, fmt.Errorf("Price must be > 0, is currently: %v. Refusing to place order", float64(purchasePrice)/1000)
	}

	order := NewBuyOrder(userID, ticker, quantity, purchasePrice, time.Now())
	err = s.buyOrders.Insert(order)
	if err != nil {
		return nil, err
	}

	err = s.usersService.UpdateForBuyOrderByUserID(order.GetUserID(), ticker, quantity, purchasePrice)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *_Service) CreateSellOrder(userID, ticker string, quantity int) (SellOrder, error) {
	price, err := stockPrice(ticker)
	if err != nil {
		return nil, err
	}

	if price <= 0 {
		return nil, fmt.Errorf("Price must be > 0, is currently: %v. Refusing to place order", float64(price)/1000)
	}

	err = s.usersService.UpdateForSellOrderByUserID(userID, ticker, quantity, price)
	if err != nil {
		return nil, err
	}

	order := NewSellOrder(userID, ticker, quantity, price, time.Now())
	err = s.sellOrders.Insert(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *_Service) GetBuyOrder(userID, id string) (BuyOrder, error) {
	_order := _BuyOrder{}

	err := s.buyOrders.Find(bson.M{"user_id": userID, "id": id}).One(&_order)
	if err != nil {
		return nil, err
	}

	return &_order, nil
}

func (s *_Service) GetSellOrder(userID, id string) (SellOrder, error) {
	_order := _SellOrder{}

	err := s.sellOrders.Find(bson.M{"user_id": userID, "id": id}).One(&_order)
	if err != nil {
		return nil, err
	}

	return &_order, nil
}

func (s *_Service) ListBuyOrders(userID string) ([]BuyOrder, error) {
	var _orders []*_BuyOrder

	err := s.buyOrders.Find(bson.M{"user_id": userID}).All(&_orders)
	if err != nil {
		return nil, err
	}

	orders := make([]BuyOrder, len(_orders))
	for i, _order := range _orders {
		orders[i] = _order
	}

	return orders, nil
}

func (s *_Service) ListSellOrders(userID string) ([]SellOrder, error) {
	var _orders []*_SellOrder

	err := s.sellOrders.Find(bson.M{"user_id": userID}).All(&_orders)
	if err != nil {
		return nil, err
	}

	orders := make([]SellOrder, len(_orders))
	for i, _order := range _orders {
		orders[i] = _order
	}

	return orders, nil
}

func stripCommas(str string) string {
  return strings.Replace(str, ",", "", -1)
}

func stockPrice(ticker string) (int, error) {
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v10/finance/quoteSummary/%v?modules=summaryDetail", ticker)
	response, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	if response.StatusCode != 200 {
		return 0, fmt.Errorf("Non 200 status code received from finance.yahoo.com: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	stockResponse := struct {
		QuoteSummary struct {
			Result []struct {
				SummaryDetail struct {
					Bid struct {
						Fmt string `json:"fmt"`
					} `json:"bid"`
				} `json:"summaryDetail"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}{}

	err = json.Unmarshal(data, &stockResponse)
	if err != nil {
		return 0, err
	}

	if len(stockResponse.QuoteSummary.Result) < 1 {
		return 0, fmt.Errorf("Received less than one result")
	}

	priceStr := stockResponse.QuoteSummary.Result[0].SummaryDetail.Bid.Fmt
	price, err := strconv.ParseFloat(stripCommas(priceStr), 64)
	if err != nil {
		return 0, err
	}

	return int(price * 1000), nil
}
