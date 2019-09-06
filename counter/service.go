package main

import (
	"time"
)

// Transaction: For 2PC
type Transaction struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []Item    `json:"data"`
}

// ItemService: Store transactions and item data
type ItemService struct {
	Items        []Item
	Transactions []Transaction
}

// TransAppend: Append new transaction
func (s *ItemService) TransAppend(t Transaction) {
	s.Transactions = append(s.Transactions, t)
}

// DoTrans: Transfor transaction to items
// And remove transaction after save items
func (s *ItemService) DoTrans(id string) {
	for i, t := range s.Transactions {
		if t.ID != id {
			continue
		}

		s.Items = append(s.Items, t.Data...)
		s.Transactions = append(s.Transactions[:i], s.Transactions[i+1:]...)
		break
	}
}

// AbortTrans: Remove certain transaction
func (s *ItemService) AbortTrans(id string) {
	for i, t := range s.Transactions {
		if t.ID != id {
			continue
		}

		s.Transactions = append(s.Transactions[:i], s.Transactions[i+1:]...)
		break
	}
}

// Item: Tenant data
type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

// Count: Tenant count
type Count struct {
	Count int `json:"count"`
}
