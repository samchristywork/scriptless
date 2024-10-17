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

func createHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		age := r.FormValue("age")

		_, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
		if err != nil {
			http.Error(w, "Unable to create entry", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/read", http.StatusSeeOther)
		return
	}


	tmpl, err := template.ParseFiles("templates/base.html", "templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pageName := r.URL.Query().Get("page")

	if pageName == "" {
		pageName = "users"
	}

	if pageName == "users" {
		rows, err := db.Query("SELECT name, age FROM users;")
		if err != nil {
			http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []struct {
			Name string
			Age  string
		}

		for rows.Next() {
			var name, age string
			if err := rows.Scan(&name, &age); err != nil {
				http.Error(w, "Unable to read data", http.StatusInternalServerError)
				return
			}
			data = append(data, struct{ Name, Age string }{name, age})
		}

		tmpl, err := template.ParseFiles("templates/base.html", "templates/read.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, struct {
			Header string
			Data   []struct {
				Name string
				Age  string
			}
		}{
			Header: "Data List",
			Data:   data,
		})

		return
	}

	pageNotFoundHandler(w)
}
