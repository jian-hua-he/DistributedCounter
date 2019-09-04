package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newTestServer() *httptest.Server {
	service := ItemService{
		Items: map[string]Item{},
	}
	mux := http.NewServeMux()
	mux.Handle("/items/", &getItemHandler{
		ItemService: &service,
	})
	mux.Handle("/items", &postItemHandler{
		ItemService: &service,
	})
	return httptest.NewServer(mux)
}

func sendRequest(method string, url string, buf *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error

	// Pass buf directly will cause nil pointer error if buf is nil
	if buf == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, buf)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func TestPostNoItem(t *testing.T) {
	// Start a new server
	server := newTestServer()
	defer server.Close()

	// Sending get request
	url := fmt.Sprintf("%s/items/Foo/count", server.URL)
	resp, err := sendRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error occurred in sending get request: error %#v", err.Error())
	}
	defer resp.Body.Close()

	// Assertion
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get handler returned unexpected status code: got %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	var count Count
	if err := json.NewDecoder(resp.Body).Decode(&count); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Get handler returned unexpected body: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	expect := Count{Count: 0}
	if !reflect.DeepEqual(count, expect) {
		t.Fatalf("Get handler returned unexpected body: got %#v, want %#v", count, expect)
	}
}

func TestPostOneItem(t *testing.T) {
	// Start a new server
	server := newTestServer()
	defer server.Close()

	// Prepare posted data
	item := Item{
		ID:     "1",
		Tenant: "Foo",
	}
	jsonStr, _ := json.Marshal(item)

	// Sending post request
	url := fmt.Sprintf("%s/items", server.URL)
	resp, err := sendRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Error occurred in sending post request: error %#v", err.Error())
	}
	defer resp.Body.Close()

	// Assertion
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Handler returned unexpected status code: got %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	var target Item
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Handler returned unexpected body: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(target, item) {
		t.Fatalf("Handler returned unexpected body: got %#v, want %#v", target, item)
	}

	// Sending get request
	url = fmt.Sprintf("%s/items/Foo/count", server.URL)
	resp, err = sendRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error occurred in sending get request: error %#v", err.Error())
	}
	defer resp.Body.Close()

	// Assertion
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get handler returned unexpected status code: got %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	var count Count
	if err := json.NewDecoder(resp.Body).Decode(&count); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Get handler returned unexpected body: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	expect := Count{Count: 1}
	if !reflect.DeepEqual(count, expect) {
		t.Fatalf("Get handler returned unexpected body: got %#v, want %#v", count, expect)
	}
}
