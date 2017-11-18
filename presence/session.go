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
	"sync"

	"golang.org/x/oauth2"
)

type session struct {
	id        string
	state     string
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type store struct {
	mu     sync.RWMutex
	values map[string]session
}

var sessions = store{
	values: make(map[string]session),
}

// randStr
func randStr(length int) string {
	raw := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		log.Fatalf("randStr: read failed: %v\n", err)
	}
	str := base64.RawURLEncoding.EncodeToString(raw)
	return str[:length]
}

func (store *store) create() session {
	return store.set(
		session{
			id:    randStr(30),
			state: randStr(30),
		})
}
func (store *store) set(s session) session {
	store.mu.Lock()
	store.values[s.id] = s
	store.mu.Unlock()
	return s
}
func (store *store) delete(id string) {
	store.mu.Lock()
	delete(store.values, id)
	store.mu.Unlock()
}
func (store *store) get(id string) (session, bool) {
	store.mu.RLock()
	s, ok := store.values[id]
	store.mu.RUnlock()
	return s, ok
}

// getSession
func getSession(r *http.Request) (session, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("session: get cookie failed: %v\n", err)
		return sessions.create(), false
	}

	s, ok := sessions.get(cookie.Value)
	if !ok {
		log.Printf("session: not found\n")
		return sessions.create(), false
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
	sessions.delete(s.id)
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
	sessions.set(s)
	fmt.Printf("session: logged as %+v\n", s)

	return nil
}
