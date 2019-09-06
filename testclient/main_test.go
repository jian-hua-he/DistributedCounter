package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

const (
	COORDINATOR_HOST = "coordinator"
	COUNTER_HOST     = "counter"
)

func TestHealth(t *testing.T) {
	url := fmt.Sprintf("http://%s/health", COUNTER_HOST)
	resp, err := GET(url)
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		t.Errorf("TestHealth failed: Error %#v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestHealth failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}
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
