package usercontext

import "context"

// User is the type of value stored in the Contexts.
type User struct {
	ID string
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in context. If there
// is no user in the context, this function will crash
func FromContext(ctx context.Context) *User {
	return ctx.Value(userKey).(*User)
}
