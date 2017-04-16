package usersservice

import mgo "gopkg.in/mgo.v2"

// UsersService manages CRUD for buy & sell users
type UsersService interface {
	CreateUser() error
}

// New constructs a new UsersService that will
// persist data using the provided mongo session
func New(mongoDB *mgo.Session) UsersService {
	return nil
}
