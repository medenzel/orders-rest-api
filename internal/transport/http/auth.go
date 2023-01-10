package http

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

// validateToken - validates an incoming jwt token
func validateToken(accessToken string) bool {
	mySigningKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("An unauthorized request has been made")
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("Authorization header could not be parsed")
			return
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("Could not validate token")
			return
		}
	}
}
