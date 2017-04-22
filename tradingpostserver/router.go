package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/royvandewater/trading-post/auth0creds"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/profilescontroller"
	"github.com/royvandewater/trading-post/sellorderscontroller"
	"github.com/royvandewater/trading-post/userscontroller"
	"github.com/royvandewater/trading-post/usersservice"
	"github.com/urfave/negroni"
)

func newRouter(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) http.Handler {
	usersService := usersservice.New(auth0Creds, mongoDB)
	ordersService := ordersservice.New(mongoDB, usersService)

	buyOrdersController := buyorderscontroller.New(ordersService)
	profilesController := profilescontroller.New(usersService)
	sellOrdersController := sellorderscontroller.New(ordersService)
	usersController := userscontroller.New(usersService)

	profileRouter := mux.NewRouter().PathPrefix("/profile").Subrouter()
	profileRouter.Methods("GET").Path("/").HandlerFunc(profilesController.Get)

	profileRouter.Methods("POST").Path("/buy-orders").HandlerFunc(buyOrdersController.Create)
	profileRouter.Methods("GET").Path("/buy-orders").HandlerFunc(buyOrdersController.List)
	profileRouter.Methods("GET").Path("/buy-orders/{id}").HandlerFunc(buyOrdersController.Get)

	profileRouter.Methods("POST").Path("/sell-orders").HandlerFunc(sellOrdersController.Create)

	router := mux.NewRouter()
	router.Methods("GET").Path("/callback").HandlerFunc(usersController.Login)
	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(usersController.Authenticate),
		negroni.Wrap(profileRouter),
	))
	router.Methods("GET").Handler(http.FileServer(http.Dir("html/")))

	return router
}
