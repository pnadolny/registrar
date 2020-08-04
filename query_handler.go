package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
