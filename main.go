package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
	Age      int    `json:"age"`
	Height   int    `json:"height"`
}

// Make a array of users to manage (simulate a db layer, but in memory)
var users = make([]user, 0)

// a mutual exclusion lock to use when adding to the users array. Probably
// overkill for this sample app, but why not.
var mutex = &sync.Mutex{}

func main() {
	// allocate a new server mux a pure Go server.
	mux := http.NewServeMux()
	mux.Handle("/register", http.HandlerFunc(handleRegister))
	mux.Handle("/save", http.HandlerFunc(handleSave))
	mux.Handle("/query", http.HandlerFunc(handleQuery))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// A helper method to unmarshall the JSON into a user
func readUser(r *http.Request) user {
	var u user
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &u)
	return u
}
