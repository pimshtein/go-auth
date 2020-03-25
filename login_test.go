package main

import (
	"bytes"
	"database/sql"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	// check if status 200
	body := []byte(testLoginData)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, login, http.StatusOK, "token")

	// check if status 401
	body = []byte(`{"login":"test", "password":"testTest"}`)
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, login, http.StatusUnauthorized, "error")

	// check if status 400 (if json is invalid, f.e. log instead login)
	// Need to apply json schema validation because unmarshal in the missing tags (for example, log instead of login), puts an empty value.

	// check if status 404
	body = []byte(`{"login":"p", "password":"test"}`)
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, login, http.StatusNotFound, "error")

	// check if status 400 (if body is wrong)
	body = []byte(``)
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, login, http.StatusBadRequest, "error")
}

func TestGetLoginData(t *testing.T) {
	// Check if no errors
	body := []byte(testLoginData)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if nil != err {
		t.Fatal(err)
	}
	errorData, loginData := getLoginData(req)
	if nil != errorData {
		t.Fatal(errorData)
	}
	if "test" != loginData.Login || "test" != loginData.Password {
		t.Fatal("Error on check struct of loginData")
	}
	// Check if bad request body
	body = nil
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if nil != err {
		t.Fatal(err)
	}
	errorData, loginData = getLoginData(req)
	if nil == errorData {
		t.Fatal(errorData)
	}
	if http.StatusBadRequest != errorData.Code {
		t.Fatal("Error on check bad request in getLoginData")
	}
	// Check if bad unmarshal json body
	body = []byte(`{"login":"test", "password":"test"`)
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if nil != err {
		t.Fatal(err)
	}
	errorData, loginData = getLoginData(req)
	if nil == errorData {
		t.Fatal(errorData)
	}
	if http.StatusBadRequest != errorData.Code {
		t.Fatal("Error on check bad request in getLoginData")
	}
}

func TestGetAccount(t *testing.T) {
	// Check if account is exist in DB
	loginData := LoginData{"test", "test"}
	errorData, _ := getAccount(&loginData)
	if nil != errorData {
		t.Fatal(errorData)
	}
	// Check if user not found by login
	loginData = LoginData{"fakeUserForTests", "test"}
	errorData, _ = getAccount(&loginData)
	if nil == errorData {
		t.Fatal(errorData)
	}
	if "User not found" != errorData.Error || http.StatusNotFound != errorData.Code {
		t.Fatal("Check if user not found by login in getAccount is fail")
	}

	// Check if error on connect to DB
	loginData = LoginData{"test", "test"}
	db.Close()
	errorData, _ = getAccount(&loginData)
	if nil == errorData {
		t.Fatal(errorData)
	}
	if !strings.Contains(errorData.Error, "Something wrong on connecting to database") ||
		http.StatusInternalServerError != errorData.Code {
		t.Fatal("Check if something wrong on connect to database in getAccount is fail")
	}
	// After close, connect to DB again
	var err error
	db, err = sql.Open("postgres", settings.DatabaseConnect)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetGlobalParameters(t *testing.T) {
	globalParams := getGlobalParameters()
	if 0 == globalParams.LoyaltyCoefficient ||
		0 == globalParams.CurrencyRate {
		t.Fail()
	}
}

func TestGenerateToken(t *testing.T) {
	// Check if token was generated
	customClaims := CustomClaims{1, 1.05, 2}
	errorData, tokenString := generateToken(customClaims)
	if 0 == len(tokenString) {
		t.Fail()
	}
	if nil != errorData {
		t.Fail()
	}
}

func TestSaveTokenStorage(t *testing.T) {
	tokenString := "test"
	errorData := saveTokenToStorage(tokenString)
	if nil != errorData {
		t.Fatal(errorData.Error)
	}
	// Check if was error
	httpmock.DeactivateAndReset()
	errorData = saveTokenToStorage(tokenString)
	if nil == errorData || http.StatusInternalServerError != errorData.Code {
		t.Fatal("Error data is nil")
	}
	// Mock response from token storage
	httpmock.Activate()
	httpmock.RegisterResponder(http.MethodPut, settings.TokenStorage+"/token",
		httpmock.NewStringResponder(200, ``))
}

func TestGetMd5HashString(t *testing.T) {
	password := "test"
	testMd5Password := "098f6bcd4621d373cade4e832627b4f6"
	md5Password := getMd5HashString(password)
	if testMd5Password != md5Password {
		t.Fail()
	}
}

func TestGetCustomClaims(t *testing.T) {
	account := AccountData{1, "test", []byte("test")}
	globalParams := GlobalParameter{1.02, 2.5}
	customClaims := getCustomClaims(account, globalParams)
	if 1 != customClaims.UserId ||
		1.02 != customClaims.CurrencyRate ||
		2.5 != customClaims.LoyaltyCoefficient {
		t.Fail()
	}
}
