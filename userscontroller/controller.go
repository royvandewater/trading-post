package userscontroller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/royvandewater/trading-post/usersservice"
)

// UsersController handles HTTP requests
// regarding users
type UsersController interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

// New constructs a new UsersController instance
func New(usersService usersservice.UsersService) UsersController {
	return &controller{usersService: usersService}
}

type controller struct {
	usersService usersservice.UsersService
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := usersservice.ParseUse(r.Body)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
	}

	storedUser, code, err := c.usersService.CreateUser(user)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	storedUserJSON, err := storedUser.JSON()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to generate JSON response: %v", err.Error())))
	}

	w.WriteHeader(code)
	w.Write(storedUserJSON)
}
