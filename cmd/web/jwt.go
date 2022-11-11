package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtOptions struct {
	key     string
	expires time.Duration
}

func (app *application) generateJWT(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(app.jwtOptions.expires).Unix()
	claims["user"] = userId

	tokenString, err := token.SignedString([]byte(app.jwtOptions.key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (app *application) extractToken(cookie *http.Cookie) (*jwt.Token, error) {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.jwtOptions.key), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (app *application) extractClaims(cookie *http.Cookie) (jwt.MapClaims, error) {
	token, err := app.extractToken(cookie)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
