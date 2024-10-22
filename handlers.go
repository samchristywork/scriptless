package main

import (
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

	page(w, `<form method="post" action="/login">
    <h1>Login</h1>
    <label for="username">Username:</label>
    <input type="text" id="username" name="username" required>

    <label for="password">Password:</label>
    <input type="password" id="password" name="password" required>

    <button type="submit">Login</button>
  </form>`)
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

	foo(w)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pageName := onlyLetters(r.URL.Query().Get("page"))
	sort := onlyLetters(r.URL.Query().Get("sort"))
	sortDir := onlyLetters(r.URL.Query().Get("sortdir"))

	if pageName == "" {
		pageName = "users"
	}

	if sort == "" {
		sort = "Name"
	}

	if sortDir == "" {
		sortDir = "asc"
	}

	if pageName == "users" {
		rows, err := db.Query("SELECT name, age FROM users ORDER BY "+sort+" "+sortDir+";")
		if err != nil {
			http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		data := [][]string{}

		for rows.Next() {
			var name, age string
			if err := rows.Scan(&name, &age); err != nil {
				http.Error(w, "Unable to read data", http.StatusInternalServerError)
				return
			}
			data = append(data, []string{name, age})
		}

		page(w, table([]string{"Name", "Age"}, data))
		return
	}

	pageNotFoundHandler(w, r)
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	page(w, `404`)
}
