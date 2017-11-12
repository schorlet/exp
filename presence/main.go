package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	// 3-legged OAuth2 flow
	conf = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:8000/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
)

const loginPage = `<html>
	<body>Log in with <a href="/login">GitHub</a></body>
</html>`

const indexPage = `<html>
	<body>Logged as %q</body>
</html>`

// index
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, ok := getSession(r)
	if ok {
		fmt.Fprintf(w, indexPage, session.Login)
		return
	}
	fmt.Fprintf(w, loginPage)
}

// login
func login(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	session.setCookie(w)

	url := conf.AuthCodeURL(session.state, oauth2.AccessTypeOnline)
	fmt.Printf("login: url: %s\n", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callback
func callback(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(r)
	if ok {
		log.Printf("callback: user already logged\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	state := r.FormValue("state")
	if state != session.state {
		log.Printf("callback: invalid state: %q\n", state)
		session.clear(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("callback: state: %q\n", state)

	code := r.FormValue("code")
	fmt.Printf("callback: code: %q\n", code)

	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("callback: exchange failed: %v\n", err)
		session.clear(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("callback: token: %+v\n", token)

	if err = session.fetchUser(token); err != nil {
		log.Printf("callback: fetch user failed: %v\n", err)
		session.clear(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// main
func main() {
	fmt.Printf("main: conf: %+v\n", conf)

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)

	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/favicon.png", http.NotFound)
	http.HandleFunc("/opensearch.xml", http.NotFound)

	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		log.Fatal(err)
	}
}
