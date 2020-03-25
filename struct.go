package main

import "github.com/dgrijalva/jwt-go"

type JwtData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims CustomClaims `json:"custom,omitempty"`
}

type CustomClaims struct {
	UserId             int     `json:"userId"`
	CurrencyRate       float32 `json:"currencyRate,omitempty"`
	LoyaltyCoefficient float32 `json:"loyaltyCoefficient,omitempty"`
}

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ErrorType struct {
	Error string `json:"error"`
}

type AccountData struct {
	Id    int
	Login string
	Pass  []byte
}

type GlobalParameter struct {
	CurrencyRate       float32
	LoyaltyCoefficient float32
}

type ResponseStorage struct {
	Body   string
	Status int
}

type Response struct {
	Token string `json:"token"`
}

type ErrorData struct {
	Error string
	Code  int
}
