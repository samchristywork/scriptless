package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Printf("Received login details - Username: %s, Password: %s\n", username, password)

		if checkCreds(username, password) {
			s, err := randomBase64String(64)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			expirationTime := time.Now().Add(30 * time.Minute)

			mu.Lock()
			sessionStore[s] = Session{
				Username:  username,
				ExpiresAt: expirationTime,
			}
			mu.Unlock()

			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    s,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				Expires:  expirationTime,
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/read", http.StatusSeeOther)
			return
		} else {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}