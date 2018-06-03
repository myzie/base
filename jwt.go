package base

import jwt "github.com/dgrijalva/jwt-go"

// JWTClaims are custom claims extending default ones
type JWTClaims struct {
	Name    string `json:"name,omitempty"`
	Admin   bool   `json:"admin,omitempty"`
	Context string `json:"context,omitempty"`
	jwt.StandardClaims
}
