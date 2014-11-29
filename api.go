package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/context"
)

func handlerAccessLog(handler http.Handler) http.Handler {
	logHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s \"%s %s\"", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logHandler)
}

func jsonResponse(data interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid data in the cache. Cannot JSON encode: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func handlerApiUser(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"userLogin": context.Get(r, ContextUserLogin).(string),
		"userName": context.Get(r, ContextUserName).(string),
	}
	jsonResponse(data, w)
}

func initRouter() http.Handler {
	r := mux.NewRouter()
	// Look for the package directory
	_, filename, _, _ := runtime.Caller(1)
	rootDir := path.Dir(filename)
	// Routes
	r.HandleFunc("/_ping", handlerPing)
	r.HandleFunc("/login", handlerLogin)
	r.HandleFunc("/login/callback", handlerLoginCallback)
	r.HandleFunc("/api/user", handlerApiUser).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(path.Join(rootDir, "static")))).Methods("GET")
	// Middlewares
	handler := handlerAccessLog(r)
	handler = handlerReadCookie(handler)
	handler = context.ClearHandler(handler)
	return handler
}

func main() {
	address := ":1080"
	router := initRouter()
	log.Printf("Starting http server on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalf("Cannot start the http server: %s", err)
	}
}
