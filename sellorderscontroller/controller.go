package sellorderscontroller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// SellOrdersController handles HTTP requests
// regarding sell orders
type SellOrdersController interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

// New constructs a new SellOrdersController instance
func New() SellOrdersController {
	return &controller{}
}

type controller struct{}

func (c *controller) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(501)
}
