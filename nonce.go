package main

import (
	"github.com/speps/go-hashids"
	"time"
)

const (
	salt = "Dutch van der Linde"
)

func newNonce() (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	return h.Encode([]int{int(now.Unix())})
}
