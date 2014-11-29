package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/context"
)

var oauthConfig *oauth.Config
var cookieSecret []byte
const cookieName = "posts-auth"
type ContextKey int
const ContextUserLogin ContextKey = 0
const ContextUserName ContextKey = 1

func init() {
	// Oauth init
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		fmt.Fprintf(os.Stderr, "This app must be registered on Github.\n")
		fmt.Fprintf(os.Stderr, "If already done, please set the env vars: `GITHUB_CLIENT_ID' and `GITHUB_CLIENT_SECRET'.\n")
		os.Exit(1)
	}
	oauthConfig = &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		RedirectURL:  "http://localhost:1080/login/callback",
	}
	// Cookie init
	cookieSecret = make([]byte, 64)
	_, err := rand.Read(cookieSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot init a random secret.\n")
		os.Exit(1)
		return
	}
}

func handlerReadCookie(handler http.Handler) http.Handler {
	cookieHandler := func(w http.ResponseWriter, r *http.Request) {
		// There are some endpoints where we disable auth
		publicEndpoints := []string{"/login", "/login/callback", "/_ping"}
		for _, endpoint := range publicEndpoints {
			if r.URL.Path == endpoint {
				handler.ServeHTTP(w, r)
				return
			}
		}
		// Try to read the cookie
		if cookie, err := r.Cookie(cookieName); err == nil {
			value := make(map[string]string)
			s := securecookie.New(cookieSecret, nil)
			if err = s.Decode(cookieName, cookie.Value, &value); err == nil {
				context.Set(r, ContextUserLogin, value["login"])
				context.Set(r, ContextUserName, value["name"])
				handler.ServeHTTP(w, r)
				return
			}
		}
		// Nothing in the cookie, let's redirect to the login endpoint
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	return http.HandlerFunc(cookieHandler)
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Init full login via oauth
	http.Redirect(w, r, oauthConfig.AuthCodeURL("login"), http.StatusFound)
}

func handlerLoginCallback(w http.ResponseWriter, r *http.Request) {
	t := &oauth.Transport{Config: oauthConfig}
	t.Exchange(r.FormValue("code"))
	client := github.NewClient(t.Client())
	user, _, err := client.Users.Get("")
	if err != nil {
		w.WriteHeader(403)
		fmt.Fprintf(w, "Cannot authenticate")
		return
	}
	// Auth ok, set cookie
	value := map[string]string{
		"login": *user.Login,
		"name":  *user.Name,
	}
	s := securecookie.New(cookieSecret, nil)
	if encoded, err := s.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
	// All done, let's go back to root
	http.Redirect(w, r, "/", http.StatusFound)
}
