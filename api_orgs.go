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
	orgs := []Organization{}
	db := getDB()
	db.Find(&orgs)
	jsonResponse(w, 200, orgs)
}

func handlerApiOrgsPost(w http.ResponseWriter, r *http.Request) {
	if user := context.Get(r, ContextUser).(User); user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	input := Organization{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		jsonError(w, 400, fmt.Sprintf("Cannot read request input: %s", err))
		return
	}
	if input.Name == "" || input.Email == "" {
		jsonError(w, 400, "Name and Email are required")
		return
	}
	org := Organization{
		Name: input.Name,
		Description: input.Description,
		Url: input.Url,
		Email: input.Email,
	}
	db := getDB()
	db.NewRecord(org)
	db.Create(&org)
}

func handlerApiOrgPut(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	if user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	input := Organization{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		jsonError(w, 400, fmt.Sprintf("Cannot read request input: %s", err))
		return
	}
	vars := mux.Vars(r)
	org := Organization{}
	db := getDB()
	db.Where("id = ?", vars["id"]).First(&org)
	org.Name = input.Name
	org.Description = input.Description
	org.Url = input.Url
	org.Email = input.Email
	db.Save(&org)
}

func handlerApiOrgGet(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	if user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	vars := mux.Vars(r)
	org := Organization{}
	db := getDB()
	db.Where("id = ?", vars["id"]).First(&org)
	if org.Id <= 0 {
		jsonError(w, 404, "Cannot find org")
		return
	}
	jsonResponse(w, 200, org)
}

func handlerApiOrgDelete(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, ContextUser).(User)
	if user.SuperAdmin != true {
		jsonError(w, 403, "Unauthorized")
		return
	}
	vars := mux.Vars(r)
	org := Organization{}
	db := getDB()
	db.Where("id = ?", vars["id"]).First(&org)
	if org.Id <= 0 {
		return
	}
	db.Unscoped().Delete(&org)
}
