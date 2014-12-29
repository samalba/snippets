package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func handlerAccessLog(handler http.Handler) http.Handler {
	logHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s \"%s %s\"", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logHandler)
}

func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid data in the cache. Cannot JSON encode: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
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
	initApiRoutes(r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(path.Join(rootDir, "static")))).Methods("GET")
	// Middlewares
	handler := handlerAccessLog(r)
	handler = handlerRequireAuth(handler)
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
