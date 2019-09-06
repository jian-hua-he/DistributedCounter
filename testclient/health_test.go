package testclient

import (
	"fmt"
	"net/http"
	"testing"
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
		t.Errorf("TestHealth failed: error %#v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestHealth failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}
}
