package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/schorlet/exp/presence/session"
	oauth "github.com/schorlet/exp/presence/slack"
	// oauth "github.com/schorlet/exp/presence/github"
)

const indexPage = `<html>
	<body>Logged as %q (%s) <a href="/logout">Log out</a></body>
</html>`

const loginPage = `<html>
	<body>Log in with <a href="/login">GitHub</a></body>
</html>`

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sess, ok := session.Get(r)
	if ok {
		fmt.Fprintf(w, indexPage, sess.User.Name, sess.User.ID)
		return
	}
	fmt.Fprintf(w, loginPage)
}

func login(w http.ResponseWriter, r *http.Request) {
	sess, ok := session.Get(r)
	if ok {
		log.Println("login: user already logged")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sess.SetCookie(w)
	url := oauth.LoginURL(sess.State)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session.Clear(w, r, "user logout")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	sess, ok := session.Get(r)
	fmt.Printf("callback session: %+v\n", sess)
	if ok {
		log.Println("callback: user already logged")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	state := r.FormValue("state")
	if state != sess.State {
		session.Clear(w, r, fmt.Sprintf(
			"callback: invalid state: %q, expected: %q\n", state, sess.State))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("callback: state: %q\n", state)

	code := r.FormValue("code")
	fmt.Printf("callback: code: %q\n", code)

	token, err := oauth.Exchange(code)
	if err != nil {
		session.Clear(w, r, fmt.Sprintf("callback: exchange failed: %v\n", err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("callback: token: %+v\n", token)

	user, err := oauth.GetUser(token)
	if err != nil {
		session.Clear(w, r, fmt.Sprintf("callback: getUser failed: %v\n", err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sess.SetUser(session.User{
		ID:   user.ID,
		Name: user.Name,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/callback", callback)

	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/favicon.png", http.NotFound)
	http.HandleFunc("/opensearch.xml", http.NotFound)

	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		log.Fatal(err)
	}
}
