package main

import (
	"math/rand"
	"time"
)

const (
	// Character set for random string
	CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Time format for id prefix
	ID_PREFIX_FORMAT = "200601021504050700"
)

// RendString: Generate random string by length
func RandString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	b := make([]byte, length)
	for i := range b {
		b[i] = CHARSET[random.Intn(len(CHARSET))]
	}

	return string(b)
}

// GenID: Generate ID with datetime and random string with 5 chars
func GenID(t time.Time) string {
	return t.Format(ID_PREFIX_FORMAT) + "-" + RandString(5)
}
