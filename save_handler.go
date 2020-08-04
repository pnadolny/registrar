package main

import "net/http"

func handleSave(writer http.ResponseWriter, request *http.Request) {
	user := readUser(request)
	for i, v := range users {
		if v.Nonce == user.Nonce {
			store.update(i, user)
			return
		}
		if v.Username == user.Username && v.Password == user.Password {
			store.update(i, user)
			return
		}
	}
	writer.WriteHeader(http.StatusUnauthorized)
}
