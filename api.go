package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/users/me", handlerApiUsersMeGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersPost).Methods("POST")
	r.HandleFunc("/api/users/{login}", handlerApiUserGet).Methods("GET")
	r.HandleFunc("/api/users/{login}", handlerApiUserPut).Methods("PUT")
	r.HandleFunc("/api/users/{login}", handlerApiUserDelete).Methods("DELETE")
	r.HandleFunc("/api/orgs", handlerApiOrgsGet).Methods("GET")
	r.HandleFunc("/api/orgs", handlerApiOrgsPost).Methods("POST")
	r.HandleFunc("/api/orgs/{id}", handlerApiOrgGet).Methods("GET")
	r.HandleFunc("/api/orgs/{id}", handlerApiOrgPut).Methods("PUT")
	r.HandleFunc("/api/orgs/{id}", handlerApiOrgDelete).Methods("DELETE")
}

func jsonError(w http.ResponseWriter, code int, msg string) {
	err := map[string]string{
		"error": msg,
	}
	jsonResponse(w, code, err)
}
