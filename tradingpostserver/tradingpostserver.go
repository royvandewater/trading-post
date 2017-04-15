package tradingpostserver

import (
	"fmt"
	"net/http"

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
func New(port int, mongoDBURL string) TradingPostServer {
	return &httpServer{
		port:       port,
		mongoDBURL: mongoDBURL,
	}
}

type httpServer struct {
	port       int
	mongoDBURL string
}

func (server *httpServer) Run() error {
	mongoDB, err := mgo.Dial(server.mongoDBURL)
	if err != nil {
		return err
	}
	defer mongoDB.Close()

	addr := fmt.Sprintf(":%v", server.port)
	return http.ListenAndServe(addr, newRouter(mongoDB))
}
