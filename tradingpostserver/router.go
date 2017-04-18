package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"github.com/royvandewater/trading-post/auth0creds"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/sellorderscontroller"
	"github.com/royvandewater/trading-post/userscontroller"
	"github.com/royvandewater/trading-post/usersservice"
)

func newRouter(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) http.Handler {
	ordersService := ordersservice.New(mongoDB)
	usersService := usersservice.New(auth0Creds, mongoDB)

	buyOrdersController := buyorderscontroller.New(ordersService)
	sellOrdersController := sellorderscontroller.New()
	usersController := userscontroller.New(usersService)

	router := httprouter.New()
	router.POST("/buy-orders", buyOrdersController.Create)
	router.POST("/sell-orders", sellOrdersController.Create)
	router.GET("/callback", usersController.Login)
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.ServeFile(w, r, "html/index.html")
	})
	return router
}
