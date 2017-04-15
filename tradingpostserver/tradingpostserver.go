package tradingpostserver

import (
	"fmt"
	"net/http"
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
func New(port int) TradingPostServer {
	return &httpServer{port: port}
}

type httpServer struct {
	port int
}

func (server *httpServer) Run() error {
	addr := fmt.Sprintf(":%v", server.port)
	return http.ListenAndServe(addr, newRouter())
}
