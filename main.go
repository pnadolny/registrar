package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	salt = "Dutch van der Linde"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
	Age      int    `json:"age"`
	Height   int    `json:"height"`
}

var users = make([]user, 0)
var mutex = &sync.Mutex{}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/register", http.HandlerFunc(handleRegister))
	mux.Handle("/save", http.HandlerFunc(handleSave))
	mux.Handle("/query", http.HandlerFunc(handleQuery))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleRegister(writer http.ResponseWriter, request *http.Request) {

	u, err := readUser(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Password == "" {
		var hash, _ = newNonce()
		mutex.Lock()
		defer mutex.Unlock()

		newUser := &user{
			Username: u.Username,
			Password: "",
			Nonce:    hash,
			Age:      0,
			Height:   0,
		}
		users = append(users, *newUser)
		res, _ := json.Marshal(newUser)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write(res)
		return
	}
}

func readUser(r *http.Request) (user, error) {
	var p user
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func handleSave(writer http.ResponseWriter, request *http.Request) {
	user, _ := readUser(request)
	for k, v := range users {
		if v.Nonce == user.Nonce {
			users[k].Age = user.Age
			users[k].Height = user.Height
			users[k].Password = user.Password
			return
		}
		if v.Username == user.Username && v.Password == user.Password {
			users[k].Age = user.Age
			users[k].Height = user.Height
			return
		}
	}
	writer.WriteHeader(http.StatusUnauthorized)
}

func handleQuery(writer http.ResponseWriter, request *http.Request) {

	var readParam = func(key string, request *http.Request) int {
		keys, ok := request.URL.Query()[key]
		if ok && len(keys[0]) > 0 {
			if value, err := strconv.Atoi(keys[0]); err == nil {
				return value
			}
		}
		return -1
	}

	for _, user := range users {
		if user.Age == readParam("age", request) {
			res, _ := json.Marshal(map[string]int{"age": user.Age})
			writer.Write(res)
			return
		}
		if user.Height == readParam("height", request) {
			res, _ := json.Marshal(map[string]int{"height": user.Height})
			writer.Write(res)
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)

}
