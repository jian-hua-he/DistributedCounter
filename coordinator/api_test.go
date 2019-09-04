package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPostItems(t *testing.T) {
	server := httptest.NewServer(new(postItemHandler))
	defer server.Close()

	item := Item{
		ID:     "1",
		Tenant: "Foo",
	}
	jsonStr, _ := json.Marshal(item)
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error occurred in sending post request: error %#v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Handler returned unexpected status code: got %#v, want %#v", resp.Status, http.StatusOK)
	}

	var target Item
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("Handler returned unexpected body: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(target, item) {
		t.Fatalf("Handler returned unexpected body: got %#v, want %#v", target, item)
	}
}
