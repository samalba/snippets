package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func handlerApiOrgsGet(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	users := []User{}
	db := getDB()
	db.Find(&users)
	jsonResponse(w, 200, users)
}

func handlerApiOrgGet(w http.ResponseWriter, r *http.Request) {
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
