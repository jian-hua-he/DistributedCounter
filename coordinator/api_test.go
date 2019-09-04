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
	req, err := http.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func TestPostItems(t *testing.T) {
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
		t.Fatalf("Handler returned unexpected status code: got %#v, want %#v", resp.Status, http.StatusOK)
	}

	var target Item
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Handler returned unexpected body: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(target, item) {
		t.Fatalf("Handler returned unexpected body: got %#v, want %#v", target, item)
	}
}
