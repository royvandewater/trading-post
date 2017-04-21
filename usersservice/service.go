package usersservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/royvandewater/trading-post/auth0creds"
	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UsersService manages CRUD for buy & sell users
type UsersService interface {
	// GetProfile retrieves a profile for a userID
	GetProfile(userID string) (Profile, error)

	// Login exchanges the code for a user profile
	// and then upserts the user profile into persistent
	// storage
	Login(code string) (User, int, error)

	// SubstractRichesByUserID removes riches from the user. Will
	// return an error if the user cannot be found
	SubstractRichesByUserID(userID string, amount float32) error

	// UserIDForAccessToken verifies the RS256 signature
	// of a JWT access token
	UserIDForAccessToken(accessToken string) (string, error)
}

// New constructs a new UsersService that will
// persist data using the provided mongo session
func New(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) UsersService {
	profiles := mongoDB.DB("tradingPost").C("profiles")
	return &_Service{auth0Creds: auth0Creds, profiles: profiles}
}

type _Service struct {
	auth0Creds auth0creds.Auth0Creds
	profiles   *mgo.Collection
}

func (s *_Service) GetProfile(userID string) (Profile, error) {
	profile := &_Profile{}

	err := s.profiles.Find(bson.M{"user_id": userID}).One(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Login finds or creates a user in the database
func (s *_Service) Login(code string) (User, int, error) {
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

	userID := profile.UserID

	_, err = s.profiles.Upsert(bson.M{"user_id": userID}, &profile)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = s.profiles.Update(bson.M{"user_id": userID, "riches": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"riches": 0}})
	if err != nil && err != mgo.ErrNotFound {
		return nil, http.StatusInternalServerError, err
	}

	err = s.profiles.Find(bson.M{"user_id": userID}).One(&profile)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user := _User{
		IDToken:     token.Extra("id_token").(string),
		AccessToken: token.AccessToken,
		Profile:     profile,
	}
	return &user, 0, nil
}

func (s *_Service) SubstractRichesByUserID(userID string, amount float32) error {
	query := bson.M{"user_id": userID}
	update := bson.M{"$inc": bson.M{"riches": -1 * amount}}

	return s.profiles.Update(query, update)
}

func (s *_Service) UserIDForAccessToken(accessToken string) (string, error) {
	claims := &_ProfileClaims{}

	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.auth0Creds.ClientSecret), nil
	})

	if err != nil {
		return "", err
	}

	return claims.UserID, nil
	// return "", nil
}

type _ProfileClaims struct {
	jwt.StandardClaims

	UserID string `json:"sub"`
}
