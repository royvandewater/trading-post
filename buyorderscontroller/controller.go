package buyorderscontroller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// BuyOrdersController handles HTTP requests
// regarding buy orders
type BuyOrdersController interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

// New constructs a new BuyOrdersController instance
func New() BuyOrdersController {
	return &controller{}
}

type controller struct{}

func (c *controller) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(501)
}
