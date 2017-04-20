package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/royvandewater/trading-post/auth0creds"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/sellorderscontroller"
	"github.com/royvandewater/trading-post/userscontroller"
	"github.com/royvandewater/trading-post/usersservice"
	"github.com/urfave/negroni"
)

func newRouter(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) http.Handler {
	ordersService := ordersservice.New(mongoDB)
	usersService := usersservice.New(auth0Creds, mongoDB)

	buyOrdersController := buyorderscontroller.New(ordersService)
	sellOrdersController := sellorderscontroller.New()
	usersController := userscontroller.New(usersService)

	router := mux.NewRouter()
	router.Methods("GET").Path("/callback").HandlerFunc(usersController.Login)
	router.Methods("GET").Handler(http.FileServer(http.Dir("html/")))

	profileRouter := mux.NewRouter().PathPrefix("/profile").Subrouter()
	profileRouter.Methods("POST").Path("/buy-orders").HandlerFunc(buyOrdersController.Create)
	profileRouter.Methods("POST").Path("/sell-orders").HandlerFunc(sellOrdersController.Create)
	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(usersController.Authenticate),
		negroni.Wrap(profileRouter),
	))

	return router
}
