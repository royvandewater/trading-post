package userscontroller

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/usersservice"
)

// UsersController handles HTTP requests
// regarding users
type UsersController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

// New constructs a new UsersController instance
func New(usersService usersservice.UsersService) UsersController {
	return &_Controller{usersService: usersService}
}

type _Controller struct {
	usersService usersservice.UsersService
}

func (c *_Controller) Login(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, statusCode, err := c.usersService.Login(code)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	userJSON, err := user.JSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate JSON response: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(userJSON)
}
