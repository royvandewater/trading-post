package sellorderscontroller

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/usercontext"
)

// Controller handles HTTP requests
// regarding sell orders
type Controller interface {
	// Create creates a new order and subtracts the market rate
	// from the profile's riches
	Create(rw http.ResponseWriter, r *http.Request)
	//
	// // Get retrieves a single order by id
	// Get(rw http.ResponseWriter, r *http.Request)
	//
	// // List retrieves all orders for a profile
	// List(rw http.ResponseWriter, r *http.Request)
}

// New constructs a new Controller instance
func New(ordersService ordersservice.OrdersService) Controller {
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

	order, err := controller.ordersService.CreateSellOrder(user.ID, createBody.Ticker)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	response, err := formatCreateResponse(order)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(201)
	rw.Write(response)
}
