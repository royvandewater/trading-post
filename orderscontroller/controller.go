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
	List(w http.ResponseWriter, r *http.Request)
}

// New constructs a new OrdersController instance
func New(ordersService ordersservice.OrdersService) OrdersController {
	return &_Controller{ordersService: ordersService}
}

type _Controller struct {
	ordersService ordersservice.OrdersService
}

func (controller *_Controller) Create(rw http.ResponseWriter, r *http.Request) {
	user := usercontext.FromContext(r.Context())

	createBody, err := parseCreateBody(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 422)
		return
	}

	order, code, err := controller.ordersService.Create(user.ID, createBody.Ticker)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	orderResponse, err := formatCreateResponse(order)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(201)
	rw.Write(orderResponse)
}

func (controller *_Controller) List(w http.ResponseWriter, r *http.Request) {
	// user := usercontext.FromContext(r.Context())
	//
	// ordersservice.List
}
