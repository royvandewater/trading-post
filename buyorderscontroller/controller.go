package buyorderscontroller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/royvandewater/trading-post/ordersservice"
)

// BuyOrdersController handles HTTP requests
// regarding buy orders
type BuyOrdersController interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

// New constructs a new BuyOrdersController instance
func New(ordersService ordersservice.OrdersService) BuyOrdersController {
	return &controller{ordersService: ordersService}
}

type controller struct {
	ordersService ordersservice.OrdersService
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	buyOrder, err := ordersservice.ParseBuyOrder(r.Body)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
	}

	storedBuyOrder, code, err := c.ordersService.CreateBuyOrder(buyOrder)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	storedBuyOrderJSON, err := storedBuyOrder.JSON()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to generate JSON response: %v", err.Error())))
	}

	w.WriteHeader(code)
	w.Write(storedBuyOrderJSON)
}
