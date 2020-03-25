package main

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"testing"
	"time"
)

func TestCheckJwtToken(t *testing.T) {
	// check if not error
	customClaims := CustomClaims{UserId: 1, CurrencyRate: 1.02, LoyaltyCoefficient: 5}
	claims := JwtData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
		CustomClaims: customClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		t.Fatal(err)
	}
	err = checkJwtToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	// check if error
	tokenString, err = token.SignedString([]byte(""))
	if err != nil {
		t.Fatal(err)
	}
	err = checkJwtToken(tokenString)
	if err == nil {
		t.Errorf("requestStorage returned unexpected status: got not error want error")
	}
}

func TestParseToken(t *testing.T) {
	// check if headers token is valid
	req, err := http.NewRequest("POST", "/logout", nil)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzU4NzE5NzcsImN1c3RvbSI6eyJ1c2VySWQiOiIxIn19.5aR0A06kMs2bvRM0_V9Thqubb9W3Bat3yP6YdOgn9-o"
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}

	_, err = parseToken(req)
	if err != nil {
		t.Errorf("parseToken return error on valid token: %v", token)
	}

	req.Header.Set("Authorization", "Bearer ")
	_, err = parseToken(req)
	if err == nil {
		t.Errorf("parseToken didn't return error on invalid token: %v", token)
	}

	req.Header.Set("Authorization", "")
	_, err = parseToken(req)
	if err == nil {
		t.Errorf("parseToken didn't return error on invalid token: %v", token)
	}
}
