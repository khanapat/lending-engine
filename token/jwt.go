package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func GenerateJWTToken(secret string, id int) (string, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString("jwt.issuer"),
			ExpiresAt: time.Now().Add(viper.GetDuration("jwt.expired-at")).Unix(),
		},
		AccountID: id,
	}
	tokenConfig := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenConfig.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
