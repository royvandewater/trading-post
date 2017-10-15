package profilescontroller

import (
	"fmt"
	"net/http"

	"github.com/royvandewater/trading-post/usercontext"
	"github.com/royvandewater/trading-post/usersservice"
)

// ProfilesController handles HTTP requests
// regarding users
type ProfilesController interface {
	// Get retrieves a user profile from persistent storage and returns it
	// as JSON
	Get(rw http.ResponseWriter, r *http.Request)
}

// New constructs a new ProfilesController instance
func New(usersService usersservice.UsersService) ProfilesController {
	return &_Controller{usersService: usersService}
}

type _Controller struct {
	usersService usersservice.UsersService
}

func (c *_Controller) Get(rw http.ResponseWriter, r *http.Request) {
	user := usercontext.FromContext(r.Context())

	profile, err := c.usersService.GetProfile(user.ID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Could retrieve profile: %v", err.Error()), 500)
		return
	}

	response, err := formatGetResponse(profile)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Could not generate response: %v", err.Error()), 500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
