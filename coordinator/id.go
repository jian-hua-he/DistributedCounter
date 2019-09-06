package main

import (
	"math/rand"
	"time"
)

const (
	CHARSET           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ID_PREPEND_FORMAT = "200601021504050700"
)

func RandString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	b := make([]byte, length)
	for i := range b {
		b[i] = CHARSET[random.Intn(len(CHARSET))]
	}

	return string(b)
}

func GenID(t time.Time) string {
	return t.Format(ID_PREPEND_FORMAT) + "-" + RandString(5)
}
