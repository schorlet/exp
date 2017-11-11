package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
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

// save
func (s session) save() {
	sessions[s.id] = s
}
