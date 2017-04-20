package buyorderscontroller

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/usercontext"
)

// BuyOrdersController handles HTTP requests
// regarding buy orders
type BuyOrdersController interface {
	Create(w http.ResponseWriter, r *http.Request)
}

// New constructs a new BuyOrdersController instance
func New(ordersService ordersservice.OrdersService) BuyOrdersController {
	return &controller{ordersService: ordersService}
}

type controller struct {
	ordersService ordersservice.OrdersService
}

func (c *controller) Create(rw http.ResponseWriter, r *http.Request) {
	user, ok := usercontext.FromContext(r.Context())
	if !ok {
		http.Error(rw, "Somehow got to an authenticated area without authentication", 500)
		return
	}

	buyOrder, err := ordersservice.ParseBuyOrderForUserID(user.ID, r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 422)
		return
	}

	storedBuyOrder, code, err := c.ordersService.CreateBuyOrder(buyOrder)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	storedBuyOrderJSON, err := storedBuyOrder.JSON()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate JSON response: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(code)
	rw.Write(storedBuyOrderJSON)
}
