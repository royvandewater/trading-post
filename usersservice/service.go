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

	// UpdateForBuyOrderByUserID removes riches from the user and adds
	// to the stock quantity for the given ticker. Will
	// return an error if the user cannot be found
	UpdateForBuyOrderByUserID(userID, ticker string, quantity int, price int) error

	// UpdateForSellOrderByUserID adds to the user's riches. Will
	// return an error if the user cannot be found
	UpdateForSellOrderByUserID(userID, ticker string, quantity int, amount int) error

	// UserIDForAccessToken verifies the RS256 signature
	// of a JWT access token
	UserIDForAccessToken(accessToken string) (string, error)
}

// New constructs a new UsersService that will
// persist data using the provided mongo session
func New(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) UsersService {
	profiles := mongoDB.DB("trading_post").C("profiles")

	conf := &oauth2.Config{
		ClientID:     auth0Creds.ClientID,
		ClientSecret: auth0Creds.ClientSecret,
		RedirectURL:  auth0Creds.CallbackURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%v/authorize", auth0Creds.Domain),
			TokenURL: fmt.Sprintf("https://%v/oauth/token", auth0Creds.Domain),
		},
	}

	return &_Service{auth0Creds: auth0Creds, conf: conf, profiles: profiles}
}

type _Service struct {
	auth0Creds auth0creds.Auth0Creds
	conf       *oauth2.Config
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
	token, err := s.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	profile, err := s.findOrCreateProfile(token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user := _User{
		IDToken:      token.Extra("id_token").(string),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Profile:      profile,
	}
	return &user, 200, nil
}

func (s *_Service) UpdateForBuyOrderByUserID(userID, ticker string, quantity, price int) error {
	err := s.ensureTickerIsPresent(userID, ticker)
	if err != nil {
		return err
	}

	query := bson.M{"user_id": userID, "stocks.ticker": ticker}
	update := bson.M{"$inc": bson.M{"riches": -1 * quantity * price, "stocks.$.quantity": quantity}}

	return s.profiles.Update(query, update)
}

func (s *_Service) UpdateForSellOrderByUserID(userID, ticker string, quantity, amount int) error {
	query := bson.M{
		"user_id": userID,
		"stocks": bson.M{
			"$elemMatch": bson.M{
				"ticker": ticker,
				"quantity": bson.M{
					"$gte": quantity,
				},
			},
		},
	}
	update := bson.M{"$inc": bson.M{"riches": quantity * amount, "stocks.$.quantity": -1 * quantity}}

	err := s.profiles.Update(query, update)
	if err != nil && err == mgo.ErrNotFound {
		return fmt.Errorf("Insufficient quantity available for: %v", ticker)
	}

	return err
}

func (s *_Service) UserIDForAccessToken(accessToken string) (string, error) {
	claims := &_ProfileClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.auth0Creds.ClientSecret), nil
	})

	if token.Valid {
		return claims.UserID, nil
	}
	if validationErr, ok := err.(*jwt.ValidationError); ok {
		// if jwt.ValidationErrorExpired is the only error, return
		if validationErr.Errors&^jwt.ValidationErrorExpired == 0 {
			return claims.UserID, nil
		}
		return "", err
	}

	return "", err

	// token := &oauth2.Token{AccessToken: accessToken}
	// _, err = s.findOrCreateProfile(token)
	// if err != nil {
	// 	return "", err
	// }
	//
}

func (s *_Service) ensureTickerIsPresent(userID, ticker string) error {
	query := bson.M{"user_id": userID, "stocks.ticker": bson.M{"$ne": ticker}}
	update := bson.M{"$addToSet": bson.M{"stocks": bson.M{"ticker": ticker, "quantity": 0}}}

	err := s.profiles.Update(query, update)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func (s *_Service) findOrCreateProfile(token *oauth2.Token) (*_Profile, error) {
	profile, err := s.getProfileForToken(token)
	if err != nil {
		return profile, err
	}

	err = s.profiles.Find(bson.M{"user_id": profile.UserID}).One(&profile)
	if err != nil && err != mgo.ErrNotFound {
		return profile, err
	}
	if err == nil {
		return profile, nil
	}

	_, err = s.profiles.Upsert(bson.M{"user_id": profile.UserID}, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *_Service) getProfileForToken(token *oauth2.Token) (*_Profile, error) {
	var profile _Profile

	client := s.conf.Client(oauth2.NoContext, token)
	resp, err := client.Get(fmt.Sprintf("https://%v/userinfo", s.auth0Creds.Domain))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(raw, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

type _ProfileClaims struct {
	jwt.StandardClaims

	UserID string `json:"sub"`
}
