package main

import (
	"reflect"
	"testing"
	"time"
)

func TestDoTrans(t *testing.T) {
	s := ItemService{
		Items:        []Item{},
		Transactions: []Transaction{},
	}

	items1 := []Item{
		{
			ID:     "id1",
			Tenant: "tenant1",
		},
		{
			ID:     "id2",
			Tenant: "tenant1",
		},
		{
			ID:     "id3",
			Tenant: "tenant1",
		},
	}
	tran1 := Transaction{
		ID:        "id1",
		Timestamp: time.Now(),
		Data:      items1,
	}

	items2 := []Item{
		{
			ID:     "id4",
			Tenant: "tenant1",
		},
	}
	tran2 := Transaction{
		ID:        "id2",
		Timestamp: time.Now(),
		Data:      items2,
	}

	s.TransAppend(tran1)
	s.TransAppend(tran2)

	if !reflect.DeepEqual(s.Transactions[0], tran1) {
		t.Errorf("Data was incorrect, got %#v, want %#v", s.Transactions[0], tran1)
	}
	if !reflect.DeepEqual(s.Transactions[1], tran2) {
		t.Errorf("Data was incorrect, got %#v, want %#v", s.Transactions[1], tran2)
	}

	s.DoTrans(tran2.ID)
	s.DoTrans(tran1.ID)

	if len(s.Transactions) > 0 {
		t.Errorf("Data was incorrect, got %#v, want %#v", len(s.Transactions), 0)
	}
	if len(s.Items) != 4 {
		t.Errorf("Data was incorrect, got %#v, want %#v", len(s.Items), 4)
	}

	want := append(items2, items1...)
	got := s.Items
	for i := 0; i < len(want); i += 1 {
		if want[i] != got[i] {
			t.Errorf("Data was incorrect, got %#v, want %#v", got[i], want[i])
		}
	}
}

func TestAbortTrans(t *testing.T) {
	s := ItemService{
		Items:        []Item{},
		Transactions: []Transaction{},
	}

	items1 := []Item{
		{
			ID:     "id1",
			Tenant: "tenant1",
		},
		{
			ID:     "id2",
			Tenant: "tenant1",
		},
		{
			ID:     "id3",
			Tenant: "tenant1",
		},
	}
	tran1 := Transaction{
		ID:        "id1",
		Timestamp: time.Now(),
		Data:      items1,
	}

	items2 := []Item{
		{
			ID:     "id4",
			Tenant: "tenant1",
		},
	}
	tran2 := Transaction{
		ID:        "id2",
		Timestamp: time.Now(),
		Data:      items2,
	}

	s.TransAppend(tran1)
	s.TransAppend(tran2)

	if !reflect.DeepEqual(s.Transactions[0], tran1) {
		t.Errorf("Data was incorrect, got %#v, want %#v", s.Transactions[0], tran1)
	}
	if !reflect.DeepEqual(s.Transactions[1], tran2) {
		t.Errorf("Data was incorrect, got %#v, want %#v", s.Transactions[1], tran2)
	}

	s.AbortTrans(tran1.ID)
	s.AbortTrans(tran2.ID)

	if len(s.Transactions) > 0 {
		t.Errorf("Data was incorrect, got %#v, want %#v", len(s.Transactions), 0)
	}
	if len(s.Items) != 0 {
		t.Errorf("Data was incorrect, got %#v, want %#v", len(s.Items), 4)
	}
}
