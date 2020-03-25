package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", `application/json`)
	// Parse and check token
	token, err := parseToken(r)
	if err != nil {
		writeError(w, errorToJson(err.Error()), http.StatusBadRequest)
		return
	}
	err = checkJwtToken(token)
	if err != nil {
		writeError(w, errorToJson(err.Error()), http.StatusBadRequest)
		return
	}

	// Delete token from storage
	response, err := requestToStorage(http.MethodDelete, settings.TokenStorage+"/token", token)
	if err != nil {
		writeError(w, errorToJson(err.Error()), http.StatusBadRequest)
		return
	}
	if response.Status != http.StatusOK {
		writeError(w, errorToJson(response.Body), response.Status)
		return
	}
	// If delete token is ok, do not return anything to the body, return status 200
}
