package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/users/me", handlerApiUsersMeGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersPost).Methods("POST")
	r.HandleFunc("/api/users/{login}", handlerApiUsersPut).Methods("PUT")
}

func jsonError(w http.ResponseWriter, code int, msg string) {
	err := map[string]string{
		"error": msg,
	}
	w.WriteHeader(code)
	jsonResponse(w, err)
}

func handlerApiUsersMeGet(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	jsonResponse(w, user)
}

func handlerApiUsersGet(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	users := []User{}
	db := getDB()
	db.Find(&users)
	jsonResponse(w, users)
}

func handlerApiUsersPost(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
}
