package publickey

import (
	"crypto/rsa"
	"fmt"

	"github.com/lestrrat/go-jwx/jwk"
)

// FromDomain retrieves public key bytes for a domain
func FromDomain(domain, kid string) (*rsa.PublicKey, error) {
	set, err := jwk.FetchHTTP(fmt.Sprintf("https://%v/.well-known/jwks.json", domain))
	if err != nil {
		return nil, err
	}

	key, err := set.LookupKeyID(kid)[0].Materialize()
	if err != nil {
		return nil, err
	}

	rsaPublicKey := key.(*rsa.PublicKey)
	return rsaPublicKey, nil
}

// 	response, err := http.Get(fmt.Sprintf("https://%v/.well-known/jwks.json", domain))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if response.StatusCode != 200 {
// 		return nil, fmt.Errorf("Unexpected non 200 status code: %v", response.StatusCode)
// 	}
//
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	jwksResponse := &JWKSResponse{}
// 	err = json.Unmarshal(body, jwksResponse)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	key := jwksResponse.Keys[0]
//
// 	fmt.Println("Decode exponent: ", key.Exponent)
// 	exponentBytes, err := base64.URLEncoding.DecodeString(key.Exponent)
// 	if err != nil {
// 		return nil, err
// 	}
// 	exponent, bytesRead := binary.Varint(exponentBytes)
// 	if bytesRead <= 0 {
// 		// https://golang.org/pkg/encoding/binary/#Varint
// 		return nil, fmt.Errorf("Could not decode exponent: %v", bytesRead)
// 	}
//
// 	fmt.Println("Decode modulus: ", key.Modulus)
// 	modulusBytes, err := base64.URLEncoding.DecodeString(key.Modulus)
// 	if err != nil {
// 		return nil, err
// 	}
// 	modulus := &big.Int{}
// 	modulus.SetBytes(modulusBytes)
//
// 	rsaPublicKey := rsa.PublicKey{E: int(exponent), N: modulus}
//
// 	// https://stackoverflow.com/questions/13555085/save-and-load-crypto-rsa-privatekey-to-and-from-the-disk
// 	publicKey, err := x509.MarshalPKIXPublicKey(&rsaPublicKey)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	pubBytes := pem.EncodeToMemory(&pem.Block{
// 		Type:  "RSA PUBLIC KEY",
// 		Bytes: publicKey,
// 	})
//
// 	return pubBytes, nil
// }
//
// // Key is a single RSA key
// type Key struct {
// 	Exponent string `json:"e"`
// 	Modulus  string `json:"n"`
// }
//
// // JWKSResponse is the JSON response back from https://{domain}/.well-known/jwks.json
// type JWKSResponse struct {
// 	Keys []Key `json:"keys"`
// }
