package testclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNoneItem(t *testing.T) {
	url := fmt.Sprintf("http://%s/items/foo/count", COORDINATOR_HOST)
	resp, err := GET(url)
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		t.Errorf("TestNoneItem failed: error %#v", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestNoneItem failed: got status code %#v, want %#v", resp.StatusCode, http.StatusOK)
	}

	want := Count{Count: 0}
	var got Count
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestNoneItem failed: error %#v, body %#v", err.Error(), string(bodyBytes))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TestNoneItem failed: got %#v, want %#v", got, want)
	}
}
