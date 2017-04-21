package orderscontroller

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/usercontext"
)

// OrdersController handles HTTP requests
// regarding buy orders
type OrdersController interface {
	Create(w http.ResponseWriter, r *http.Request)
}

// New constructs a new OrdersController instance
func New(ordersService ordersservice.OrdersService) OrdersController {
	return &_Controller{ordersService: ordersService}
}

type _Controller struct {
	ordersService ordersservice.OrdersService
}

func (c *_Controller) Create(rw http.ResponseWriter, r *http.Request) {
	user, ok := usercontext.FromContext(r.Context())
	if !ok {
		http.Error(rw, "Somehow got to an authenticated area without authentication", 500)
		return
	}

	createBody, err := parseCreateBody(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 422)
		return
	}

	order, code, err := c.ordersService.CreateOrder(user.ID, createBody.Ticker)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	orderJSON, err := order.JSON()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate JSON response: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(code)
	rw.Write(orderJSON)
}
