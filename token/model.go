package token

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	jwt.StandardClaims
	AccountID int `json:"accountId" example:"1"`
}
