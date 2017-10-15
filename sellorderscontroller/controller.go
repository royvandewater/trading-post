package sellorderscontroller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/usercontext"
)

// Controller handles HTTP requests
// regarding sell orders
type Controller interface {
	// Create creates a new order and adds the market rate
	// to the profile's riches
	Create(rw http.ResponseWriter, r *http.Request)

	// Get retrieves a single order by id
	Get(rw http.ResponseWriter, r *http.Request)

	// List retrieves all orders for a profile
	List(rw http.ResponseWriter, r *http.Request)
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

	order, err := controller.ordersService.CreateSellOrder(user.ID, createBody.Ticker, createBody.Quantity)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	response, err := formatCreateResponse(order)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), 500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(201)
	rw.Write(response)
}

func (controller *_Controller) Get(rw http.ResponseWriter, r *http.Request) {
	user := usercontext.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	order, err := controller.ordersService.GetSellOrder(user.ID, id)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to retrieve order: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	response, err := formatGetResponse(order)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func (controller *_Controller) List(rw http.ResponseWriter, r *http.Request) {
	user := usercontext.FromContext(r.Context())

	orders, err := controller.ordersService.ListSellOrders(user.ID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to retrieve orders: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	response, err := formatListResponse(orders)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
