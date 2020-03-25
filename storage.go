package main

import (
	"io/ioutil"
	"net/http"
)

func requestToStorage(method string, url string, data string) (ResponseStorage, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ResponseStorage{"", 500}, err
	}
	req.Header.Set("Authorization", "Bearer "+data)
	response, err := client.Do(req)
	if err != nil {
		return ResponseStorage{"", 500}, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ResponseStorage{"", 500}, err
	}
	return ResponseStorage{string(body), response.StatusCode}, nil
}
