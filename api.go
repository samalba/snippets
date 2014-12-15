package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/users/me", handlerApiUsersMe).Methods("GET")
}

func handlerApiUsersMe(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	jsonResponse(user, w)
}
