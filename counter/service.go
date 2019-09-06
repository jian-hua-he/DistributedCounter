package main

import (
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []Item    `json:"data"`
}

type ItemService struct {
	Items        []Item
	Transactions []Transaction
}

func (s *ItemService) TransAppend(t Transaction) {
	s.Transactions = append(s.Transactions, t)
}

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

func (s *ItemService) AbortTrans(id string) {
	for i, t := range s.Transactions {
		if t.ID != id {
			continue
		}

		s.Transactions = append(s.Transactions[:i], s.Transactions[i+1:]...)
		break
	}
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}
