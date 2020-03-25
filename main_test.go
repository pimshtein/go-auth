package main

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var testLoginData = `{"login":"test", "password":"test"}`

func before(t *testing.T) {
	// Mock response from token storage
	httpmock.Activate()
	httpmock.RegisterResponder(http.MethodPut, settings.TokenStorage+"/token",
		httpmock.NewStringResponder(200, ``))
}

func checkResponse(t *testing.T, req *http.Request, handlerFunction http.HandlerFunc, expectedStatus int, expectedBody string) {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunction)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	// Check the response body is what we expect.
	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func after() {
	httpmock.DeactivateAndReset()
}

func TestMain(m *testing.M) {
	t := new(testing.T)
	before(t)
	codeExit := m.Run()
	after()
	os.Exit(codeExit)
}
