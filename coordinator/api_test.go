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

func TestPostItems(t *testing.T) {
	// Start a new server
	service := ItemService{
		Items: map[string]Item{},
	}
	server := httptest.NewServer(&postItemHandler{
		ItemService: &service,
	})
	defer server.Close()

	// Prepare posted data
	item := Item{
		ID:     "1",
		Tenant: "Foo",
	}
	jsonStr, _ := json.Marshal(item)

	// Sending post request
	url := fmt.Sprintf("%s/items", server.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
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
