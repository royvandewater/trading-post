package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/royvandewater/trading-post/auth0creds"
	"github.com/royvandewater/trading-post/orderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/userscontroller"
	"github.com/royvandewater/trading-post/usersservice"
	"github.com/urfave/negroni"
)

func newRouter(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) http.Handler {
	usersService := usersservice.New(auth0Creds, mongoDB)
	ordersService := ordersservice.New(mongoDB, usersService)

	ordersController := orderscontroller.New(ordersService)
	usersController := userscontroller.New(usersService)

	router := mux.NewRouter()
	router.Methods("GET").Path("/callback").HandlerFunc(usersController.Login)
	router.Methods("GET").Handler(http.FileServer(http.Dir("html/")))

	profileRouter := mux.NewRouter().PathPrefix("/profile").Subrouter()
	profileRouter.Methods("POST").Path("/orders").HandlerFunc(ordersController.Create)

	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(usersController.Authenticate),
		negroni.Wrap(profileRouter),
	))

	return router
}
