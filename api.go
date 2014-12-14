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
	data := map[string]string{
		"login": user.Login,
		"name": user.Name,
		"company": user.Company,
		"email": user.Email,
		"avatarURL": user.AvatarURL,
		"location": user.Location,
	}
	admin := "false";
	if user.SuperAdmin == true {
		admin = "true"
	}
	data["superAdmin"] = admin
	jsonResponse(data, w)
}
