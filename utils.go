package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

func errorToJson(err string) string {
	errCombine := ErrorType{err}
	response, _ := json.Marshal(errCombine)
	return string(response[:])
}

func writeError(w http.ResponseWriter, error string, code int) {
	log.Println(error)
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func checkJwtToken(token string) error {
	_, err := jwt.ParseWithClaims(token, &JwtData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(Secret), nil
	})

	if err != nil {
		err := "Invalid parse token: " + err.Error()
		log.Println(err)
		return errors.New(err)
	}

	return nil
}

func parseToken(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	if !(len(authArr) == 2 && authArr[1] != "") {
		err := "Authentication header is invalid: " + authToken
		log.Println(err)
		return "", errors.New(err)
	}

	return authArr[1], nil
}
