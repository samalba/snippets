package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/user", handlerApiUser).Methods("GET")
}

func handlerApiUser(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"userLogin": context.Get(r, ContextUserLogin).(string),
		"userName":  context.Get(r, ContextUserName).(string),
	}
	jsonResponse(data, w)
}
