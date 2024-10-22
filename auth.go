package main

import (
	"crypto/rand"
	"encoding/base64"
	"math"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	Username  string
	ExpiresAt time.Time
}

var (
	sessionStore = make(map[string]Session)
	mu           sync.Mutex
)

func randomBase64String(l int) (string, error) {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l], nil
}

func checkCreds(username string, password string) bool {
	if username == "foo" && password == "bar" {
		return true
	}

	return false
}

func cleanExpiredSessions() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for key, session := range sessionStore {
		if now.After(session.ExpiresAt) {
			delete(sessionStore, key)
		}
	}
}

func isSessionValid(sessionID string) bool {
	mu.Lock()
	defer mu.Unlock()

	session, exists := sessionStore[sessionID]
	if !exists {
		return false
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessionStore, sessionID)
		return false
	}

	return true
}

func checkSession(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return false
	}

	if isSessionValid(cookie.Value) {
		return true
	}

	return false
}
