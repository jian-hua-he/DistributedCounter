package testclient

import (
	"bytes"
	"net/http"
)

const (
	COORDINATOR_HOST = "coordinator"
	COUNTER_HOST     = "counter"
)

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}

// GET: Send the request with GET method
func GET(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

// POST: Send the request with POST method
func POST(url string, buf *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
