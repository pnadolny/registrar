package main

import "sync"

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

type userStore struct {
}

func NewUserStore() *userStore {
	return &userStore{}
}

func (s *userStore) store(u user) {
	mutex.Lock()
	users = append(users, u)
	mutex.Unlock()
}

func (s *userStore) update(index int, u user) {
	users[index].Age = u.Age
	users[index].Height = u.Height
	users[index].Password = u.Password
}
