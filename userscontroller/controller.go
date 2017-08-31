package userscontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/royvandewater/trading-post/usercontext"
	"github.com/royvandewater/trading-post/usersservice"
)

// UsersController handles HTTP requests
// regarding users
type UsersController interface {
	// Authenticate is middleware that will gate access to any handlers
	// down the chain. After Authenticate has been called, and it has
	// successfuly next, you can access the user object using:
	// `user := usercontext.FromContext(r.Context())`
	Authenticate(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

	// Login returns a user object, along with a nested profile. The
	// id_token should be used as the a "Authorization Bearer" token
	// for example: `curl -H 'Authorization: Bearer <id_token>' http://...`
	Login(rw http.ResponseWriter, r *http.Request)

	// Token returns a user object, along with a nested profile. The
	// refresh_token should be in the POST body with a grant_type of refresh_token
	// for example:
	// ```
	// curl \
	//   -H 'Content-Type: application/json' \
	//   -d '{"grant_type": "refresh_token", "refresh_token": "<token>"}' \
	//   http://...`
	//
	Token(rw http.ResponseWriter, r *http.Request)
}

// New constructs a new UsersController instance
func New(usersService usersservice.UsersService) UsersController {
	return &_Controller{usersService: usersService}
}

type _Controller struct {
	usersService usersservice.UsersService
}

func (c *_Controller) Authenticate(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessToken, err := parseBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(rw, err.Error(), 401)
		return
	}

	userID, err := c.usersService.UserIDForAccessToken(accessToken)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Could not verify Authorization Bearer token: %v", err.Error()), 401)
		return
	}

	next(rw, r.WithContext(usercontext.NewContext(r.Context(), &usercontext.User{ID: userID})))
}

func (c *_Controller) Login(rw http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, statusCode, err := c.usersService.Login(code)
	if err != nil {
		http.Error(rw, err.Error(), statusCode)
		return
	}

	userJSON, err := formatUser(user)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate JSON response: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(userJSON)
}

func (c *_Controller) Token(rw http.ResponseWriter, r *http.Request) {
	tokenBody, err := parseTokenBody(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	err = validateTokenBody(tokenBody)
	if err != nil {
		http.Error(rw, err.Error(), 422)
		return
	}

	user, err := c.usersService.RefreshToken(tokenBody.RefreshToken)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	userJSON, err := formatUser(user)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate JSON response: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(201)
	rw.Write(userJSON)
}

func parseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("Missing required header 'Authorization'")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("header 'Authorization' was not in the required format: 'Bearer <token>'")
	}

	tokenType := strings.ToLower(parts[0])
	accessToken := parts[1]

	if tokenType != "bearer" {
		return "", fmt.Errorf("header 'Authorization' was not the required type of: 'Bearer'")
	}

	return accessToken, nil
}
