package main

import "net/http"

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
