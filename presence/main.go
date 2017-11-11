package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	confState = randomString(20)
)

// randomString
func randomString(length int) string {
	raw := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		log.Fatalf("randomString: read failed: %v\n", err)
	}
	str := base64.RawURLEncoding.EncodeToString(raw)
	return str[:length]
}

const indexPage = `<html>
	<body>Log in with <a href="/login">GitHub</a></body>
</html>`

// index
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := fmt.Fprintf(w, indexPage); err != nil {
		log.Printf("index: write response: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// login
func login(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(confState, oauth2.AccessTypeOnline)
	fmt.Printf("login: url: %s\n", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callback
func callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != confState {
		log.Printf("callback: invalid state: %q\n", state)
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("callback: token: %+v\n", token)

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("callback: new request failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	ctx = context.Background()
	client := conf.Client(ctx, token)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("callback: get authenticated user failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	user := struct {
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}{}
	dec := json.NewDecoder(res.Body)
	if err = dec.Decode(&user); err != nil {
		log.Printf("callback: decode user failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("callback: logged in as %+v\n", user)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// main
func main() {
	fmt.Printf("main: conf: %+v\n", conf)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)
	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		log.Fatal(err)
	}
}
