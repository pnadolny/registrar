package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var store = NewUserStore()

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
