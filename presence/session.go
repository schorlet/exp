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

	"golang.org/x/oauth2"
)

type session struct {
	id        string
	state     string
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

var sessions = make(map[string]session)

// randStr
func randStr(length int) string {
	raw := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		log.Fatalf("randStr: read failed: %v\n", err)
	}
	str := base64.RawURLEncoding.EncodeToString(raw)
	return str[:length]
}

// newSession
func newSession() session {
	s := session{
		id:    randStr(30),
		state: randStr(30),
	}
	sessions[s.id] = s
	return s
}

// getSession
func getSession(r *http.Request) (session, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("session: get cookie failed: %v\n", err)
		return newSession(), false
	}
	s, ok := sessions[cookie.Value]
	if !ok {
		log.Printf("session: not found\n")
		return newSession(), false
	}
	return s, s.Login != ""
}

// setCookie
func (s session) setCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    s.id,
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// clear
func (s session) clear(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	delete(sessions, s.id)
}

// fetchUser
func (s session) fetchUser(token *oauth2.Token) error {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return fmt.Errorf("session: new request failed: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	ctx := context.Background()
	client := conf.Client(ctx, token)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("session: fetch user failed: %v", err)
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	if err = dec.Decode(&s); err != nil {
		return fmt.Errorf("session: decode user failed: %v", err)
	}
	sessions[s.id] = s
	fmt.Printf("session: logged as %+v\n", s)

	return nil
}
