package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// A helper method to unmarshall the JSON into a user
func readUser(r *http.Request) user {
	var u user
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &u)
	return u
}

func handleSave(writer http.ResponseWriter, request *http.Request) {

	user := readUser(request)
	var updateUser = func(index int) {
		users[index].Age = user.Age
		users[index].Height = user.Height
		users[index].Password = user.Password
	}
	for i, v := range users {
		if v.Nonce == user.Nonce {
			updateUser(i)
			return
		}
		if v.Username == user.Username && v.Password == user.Password {
			updateUser(i)
			return
		}
	}
	writer.WriteHeader(http.StatusUnauthorized)
}

func handleQuery(writer http.ResponseWriter, request *http.Request) {

	var query = func(key string, request *http.Request) int {
		keys, ok := request.URL.Query()[key]
		if ok && len(keys[0]) > 0 {
			if value, err := strconv.Atoi(keys[0]); err == nil {
				return value
			}
		}
		return -1
	}
	for _, user := range users {
		if user.Age == query("age", request) {
			res, _ := json.Marshal(map[string]int{"age": user.Age})
			_, _ = writer.Write(res)
			return
		}
		if user.Height == query("height", request) {
			res, _ := json.Marshal(map[string]int{"height": user.Height})
			_, _ = writer.Write(res)
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)

}
