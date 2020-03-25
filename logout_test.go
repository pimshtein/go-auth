package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogout(t *testing.T) {
	httpmock.RegisterResponder(http.MethodDelete, settings.TokenStorage+"/token",
		httpmock.NewStringResponder(200, ``))
	// Get token
	body := []byte(testLoginData)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(rr, req)
	responseToken := Response{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseToken)
	if err != nil {
		t.Fatal(err)
	}

	// check if status 200
	req, err = http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+responseToken.Token)
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, logout, http.StatusOK, "")

	// check if status 400
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	checkResponse(t, req, logout, http.StatusBadRequest, "error")
}
