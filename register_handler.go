package main

import (
	"encoding/json"
	"net/http"
)

// To register, post JSON w/o a password and you will get a nonnce back
// to use later.
func handleRegister(writer http.ResponseWriter, request *http.Request) {
	u := readUser(request)
	if u.Password == "" {
		var hash, _ = newNonce()
		newUser := &user{
			Username: u.Username,
			Password: "",
			Nonce:    hash,
			Age:      0,
			Height:   0,
		}
		mutex.Lock()
		users = append(users, *newUser)
		mutex.Unlock()
		res, _ := json.Marshal(newUser)
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write(res)
	}
}
