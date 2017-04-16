package usersservice

import mgo "gopkg.in/mgo.v2"

// UsersService manages CRUD for buy & sell users
type UsersService interface {
	CreateUser(user User) (User, int, error)
}

// New constructs a new UsersService that will
// persist data using the provided mongo session
func New(mongoDB *mgo.Session) UsersService {
	users := mongoDB.DB("tradingPost").C("users")
	return &service{users: users}
}

type service struct {
	users *mgo.Collection
}

// CreateUser takes a user and stores it in the database
func (s *service) CreateUser(user User) (User, int, error) {
	err := s.users.Insert(user)
	if err != nil {
		return nil, 500, err
	}

	return user, 201, nil
}
