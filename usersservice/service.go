package usersservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/royvandewater/trading-post/auth0creds"
	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UsersService manages CRUD for buy & sell users
type UsersService interface {
	Login(code string) (User, int, error)
}

// New constructs a new UsersService that will
// persist data using the provided mongo session
func New(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) UsersService {
	profiles := mongoDB.DB("tradingPost").C("profiles")
	return &service{auth0Creds: auth0Creds, profiles: profiles}
}

type service struct {
	auth0Creds auth0creds.Auth0Creds
	profiles   *mgo.Collection
}

// Login finds or creates a user in the database
func (s *service) Login(code string) (User, int, error) {
	conf := &oauth2.Config{
		ClientID:     s.auth0Creds.ClientID,
		ClientSecret: s.auth0Creds.ClientSecret,
		RedirectURL:  s.auth0Creds.CallbackURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%v/authorize", s.auth0Creds.Domain),
			TokenURL: fmt.Sprintf("https://%v/oauth/token", s.auth0Creds.Domain),
		},
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get(fmt.Sprintf("https://%v/userinfo", s.auth0Creds.Domain))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var profile _Profile
	if err = json.Unmarshal(raw, &profile); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	s.profiles.Upsert(bson.M{"user_id": profile.UserID}, &profile)

	user := _User{
		IDToken:     token.Extra("id_token").(string),
		AccessToken: token.AccessToken,
		Profile:     profile,
	}
	return &user, 0, nil
}
