package main

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
	"time"
)

func TestRequestToStorage(t *testing.T) {
	// check if status 200
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
	respStorage, err := requestToStorage(http.MethodPut, settings.TokenStorage+"/token", tokenString)
	if respStorage.Status != http.StatusOK {
		t.Errorf("requestStorage returned unexpected status: got %v want %v", respStorage.Status, http.StatusOK)
	}
	// check if status 400
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPut, settings.TokenStorage+"/token",
		httpmock.NewStringResponder(400, ``))
	respStorage, err = requestToStorage(http.MethodPut, settings.TokenStorage+"/token", "")
	if err != nil {
		t.Fatal(err)
	}
	if respStorage.Status != http.StatusBadRequest {
		t.Errorf("requestStorage returned unexpected status: got %v want %v", respStorage.Status, http.StatusBadRequest)
	}
	// Return initial mock
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPut, settings.TokenStorage+"/token",
		httpmock.NewStringResponder(200, ``))
}
