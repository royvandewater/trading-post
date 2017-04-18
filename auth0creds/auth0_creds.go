package auth0creds

// Auth0Creds is a struct to help store
// credentials for auth0
type Auth0Creds struct {
	CallbackURL  string
	ClientID     string
	ClientSecret string
	Domain       string
}
