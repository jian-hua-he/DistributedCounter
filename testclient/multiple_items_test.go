package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestMultipleItems(t *testing.T) {
	// Prepare data
	items := []Item{
		{
			ID:     "item1",
			Tenant: "tenant1",
		},
		{
			ID:     "item2",
			Tenant: "tenant1",
		},
		{
			ID:     "item3",
			Tenant: "tenant1",
		},
		{
			ID:     "item1",
			Tenant: "tenant1",
		},
		{
			ID:     "item1",
			Tenant: "tenant2",
		},
		{
			ID:     "item2",
			Tenant: "tenant2",
		},
		{
			ID:     "item3",
			Tenant: "tenant2",
		},
		{
			ID:     "item4",
			Tenant: "tenant2",
		},
	}
	jsonStr, _ := json.Marshal(items)

	// Post to the Coordanitor
	url := fmt.Sprintf("http://%s/items", COORDINATOR_HOST)
	resp, err := POST(url, bytes.NewBuffer(jsonStr))
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		t.Errorf("TestMultipleItems failed: error %#v", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestMultipleItems failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	// Assertions
	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestMultipleItems failed: error %#v, body %#v", err.Error(), string(bodyBytes))
	}
	got1 := status.Status
	want1 := "success"
	if got1 != want1 {
		t.Errorf("TestMultipleItems failed: got %#v, want %#v", got1, want1)
	}

	// Get tenant1 data
	url = fmt.Sprintf("http://%s/items/tenant1/count", COORDINATOR_HOST)
	resp, err = GET(url)
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		t.Errorf("TestMultipleItems failed: error %#v", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestMultipleItems failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	// Assertion
	want2 := Count{Count: 3}
	var got2 Count
	if err := json.NewDecoder(resp.Body).Decode(&got2); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestMultipleItems failed: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("TestMultipleItems failed: got %#v, want %#v", got2, want2)
	}

	// Get tenant2 data
	url = fmt.Sprintf("http://%s/items/tenant2/count", COORDINATOR_HOST)
	resp, err = GET(url)
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		t.Errorf("TestMultipleItems failed: error %#v", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestMultipleItems failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	// Assertions
	want3 := Count{Count: 4}
	var got3 Count
	if err := json.NewDecoder(resp.Body).Decode(&got3); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestMultipleItems failed: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("TestMultipleItems failed: got %#v, want %#v", got3, want3)
	}
}
