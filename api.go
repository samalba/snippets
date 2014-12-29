package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func initApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/users/me", handlerApiUsersMeGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersGet).Methods("GET")
	r.HandleFunc("/api/users", handlerApiUsersPost).Methods("POST")
	r.HandleFunc("/api/users/{login}", handlerApiUserGet).Methods("GET")
	r.HandleFunc("/api/users/{login}", handlerApiUserPut).Methods("PUT")
	r.HandleFunc("/api/users/{login}", handlerApiUserDelete).Methods("DELETE")
}

func jsonError(w http.ResponseWriter, code int, msg string) {
	err := map[string]string{
		"error": msg,
	}
	jsonResponse(w, code, err)
}

func handlerApiUsersMeGet(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	jsonResponse(w, 200, user)
}

func handlerApiUsersGet(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	users := []User{}
	db := getDB()
	db.Find(&users)
	jsonResponse(w, 200, users)
}

func handlerApiUsersPost(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	input := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		jsonError(w, 400, fmt.Sprintf("Cannot read request input: %s", err))
		return
	}
	if input.Login == "" {
		jsonError(w, 400, "Login is empty")
		return
	}
	user := User{Login: input.Login}
	db := getDB()
	db.NewRecord(user)
	db.Create(&user)
}

func handlerApiUserGet(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	vars := mux.Vars(r)
	user := User{}
	db := getDB()
	db.Where("login = ?", vars["login"]).First(&user)
	if user.Id <= 0 {
		jsonError(w, 404, "Cannot find user")
		return
	}
	jsonResponse(w, 200, user)
}

func handlerApiUserPut(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	if user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	input := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		jsonError(w, 400, fmt.Sprintf("Cannot read request input: %s", err))
		return
	}
	vars := mux.Vars(r)
	if vars["login"] == user.Login {
		jsonError(w, 400, "Cannot update yourself")
		return
	}
	u := User{}
	db := getDB()
	db.Where("login = ?", vars["login"]).First(&u)
	// Only the admin flag is updateable (the rest is delegated to github and
	// updated with each new login session).
	u.SuperAdmin = input.SuperAdmin
	db.Save(&u)
}

func handlerApiUserDelete(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	if user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	vars := mux.Vars(r)
	if vars["login"] == user.Login {
		jsonError(w, 400, "Cannot delete yourself")
		return
	}
	u := User{}
	db := getDB()
	db.Where("login = ?", vars["login"]).First(&u)
	if u.Id <= 0 {
		return
	}
	db.Unscoped().Delete(&u)
}
