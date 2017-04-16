package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/sellorderscontroller"
)

func newRouter(mongoDB *mgo.Session) http.Handler {
	ordersService := ordersservice.New(mongoDB)
	// usersService := usersservice.New(mongoDB)

	buyOrdersController := buyorderscontroller.New(ordersService)
	sellOrdersController := sellorderscontroller.New()
	// usersController := userscontroller.New(usersService)

	router := httprouter.New()
	router.POST("/buy-orders", buyOrdersController.Create)
	router.POST("/sell-orders", sellOrdersController.Create)
	// router.POST("/users", usersController.Create)
	return router
}
