package tradingpostserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/sellorderscontroller"
)

func newRouter() http.Handler {
	buyOrdersController := buyorderscontroller.New()
	sellOrdersController := sellorderscontroller.New()

	router := httprouter.New()
	router.POST("/buy-orders", buyOrdersController.Create)
	router.POST("/sell-orders", sellOrdersController.Create)
	return router
}
