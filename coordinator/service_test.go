package main

import (
	"reflect"
	"testing"
)

func TestUpdateHost(t *testing.T) {
	hs := HostService{
		Hosts: map[string]Host{},
	}

	h1 := Host{
		Name:     "Host1",
		IsNew:    false,
		Attempts: 0,
	}
	h2 := Host{
		Name:     "Host2",
		IsNew:    true,
		Attempts: 0,
	}

	hs.UpdateHost(h1)
	hs.UpdateHost(h2)

	if got := hs.Hosts[h1.Name]; !reflect.DeepEqual(got, h1) {
		t.Errorf("Data was incorrect, got %#v, want %#v", got, h1)
	}
	if got := hs.Hosts[h2.Name]; !reflect.DeepEqual(got, h2) {
		t.Errorf("Data was incorrect, got %#v, want %#v", got, h2)
	}

	h1.Attempts = 2
	hs.UpdateHost(h1)
	if got := hs.Hosts[h1.Name]; !reflect.DeepEqual(got, h1) {
		t.Errorf("Data was incorrect, got %#v, want %#v", got, h1)
	}
}

func TestDeleteHost(t *testing.T) {
	hs := HostService{
		Hosts: map[string]Host{},
	}

	h1 := Host{
		Name:     "Host1",
		IsNew:    false,
		Attempts: 0,
	}
	h2 := Host{
		Name:     "Host2",
		IsNew:    true,
		Attempts: 0,
	}
	h3 := Host{
		Name:     "Host3",
		IsNew:    false,
		Attempts: 0,
	}

	hs.UpdateHost(h1)
	hs.UpdateHost(h2)
	hs.UpdateHost(h3)

	hs.DeleteHost(h2.Name)
	if got := hs.Hosts[h1.Name]; !reflect.DeepEqual(got, h1) {
		t.Errorf("Data was incorrect, got %#v, want %#v", got, h1)
	}
	if got := hs.Hosts[h3.Name]; !reflect.DeepEqual(got, h3) {
		t.Errorf("Data was incorrect, got %#v, want %#v", got, h3)
	}
	if _, ok := hs.Hosts[h2.Name]; ok {
		t.Errorf("Data was incorrect, got %#v, want %#v", ok, false)
	}
}
