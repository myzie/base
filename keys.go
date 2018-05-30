package base

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

// LoadRSAPrivateKey reads an RSA private key from the specified file
func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// LoadRSAPublicKey reads an RSA public key from the specified file
func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, err
	}
	return key, nil
}
