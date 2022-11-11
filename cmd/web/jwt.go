package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtOptions struct {
	key     string
	expires time.Duration
}

func (app *application) generateJWT(userName string) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(app.jwtOptions.expires)
	claims["user"] = userName

	tokenString, err := token.SignedString([]byte(app.jwtOptions.key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
