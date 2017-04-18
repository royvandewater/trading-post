package tradingpostserver

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/auth0creds"

	mgo "gopkg.in/mgo.v2"
)

// TradingPostServer defines the interface for a
// trading post server
type TradingPostServer interface {
	// Run causes the server to start serving up.
	// The server will run forever unless some error
	// occurs
	Run() error
}

// New instantiates a new TradingPostServer instance
func New(auth0Creds auth0creds.Auth0Creds, mongoDBURL string, port int) TradingPostServer {
	return &httpServer{
		auth0Creds: auth0Creds,
		port:       port,
		mongoDBURL: mongoDBURL,
	}
}

type httpServer struct {
	auth0Creds auth0creds.Auth0Creds
	mongoDBURL string
	port       int
}

func (server *httpServer) Run() error {
	mongoDB, err := mgo.Dial(server.mongoDBURL)
	if err != nil {
		return err
	}
	defer mongoDB.Close()

	addr := fmt.Sprintf(":%v", server.port)
	router := newRouter(server.auth0Creds, mongoDB)
	return http.ListenAndServe(addr, router)
}
