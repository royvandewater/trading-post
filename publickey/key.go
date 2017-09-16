package publickey

import (
	"crypto/rsa"
	"fmt"

	"github.com/lestrrat/go-jwx/jwk"
)

var cache map[string]*rsa.PublicKey

// FromDomain retrieves public key bytes for a domain
func FromDomain(domain, kid string) (*rsa.PublicKey, error) {
	if cache == nil {
		cache = make(map[string]*rsa.PublicKey)
	}

	if _, ok := cache[kid]; ok {
		return cache[kid], nil
	}

	set, err := jwk.FetchHTTP(fmt.Sprintf("https://%v/.well-known/jwks.json", domain))
	if err != nil {
		return nil, err
	}

	key, err := set.LookupKeyID(kid)[0].Materialize()
	if err != nil {
		return nil, err
	}

	rsaPublicKey := key.(*rsa.PublicKey)
	cache[kid] = rsaPublicKey
	return rsaPublicKey, nil
}
