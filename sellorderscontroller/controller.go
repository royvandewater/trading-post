package sellorderscontroller

import (
	"net/http"
)

// SellOrdersController handles HTTP requests
// regarding sell orders
type SellOrdersController interface {
	Create(w http.ResponseWriter, r *http.Request)
}

// New constructs a new SellOrdersController instance
func New() SellOrdersController {
	return &controller{}
}

type controller struct{}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
}
