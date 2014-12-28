package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/users/me", handlerApiUsersMe).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsers).Methods("GET")
}

func handlerApiUsersMe(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	jsonResponse(user, w)
}

func handlerApiUsers(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		w.WriteHeader(403)
		fmt.Fprintf(w, "Unauthorized")
	}
	users := []User{}
	db := getDB()
	db.Find(&users)
	jsonResponse(users, w)
}
