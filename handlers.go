package main

import (
	"html/template"
	"log"
	"fmt"
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

func buildSidebar(selected string) template.HTML {
	type Page struct {
		name string
		slug string
	}
	pages := []Page{
		{
			name: "Create",
			slug: "create",
		},
		{
			name: "Read",
			slug: "read",
		},
		{
			name: "Foo",
			slug: "foo",
		},
		{
			name: "Bar",
			slug: "bar",
		},
	}

	sidebar := template.HTML("")
	for _, s:=range pages {
		if s.name==selected {
			sidebar+=template.HTML(fmt.Sprintf(`
				<div style="color:green"><a href="/%s">%s</a></div>
			`, s.slug, s.name))
		} else {
			sidebar+=template.HTML(fmt.Sprintf(`
				<div style="color:red"><a href="/%s">%s</a></div>
			`, s.slug, s.name))
		}
	}

	return sidebar
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

	tmpl.Execute(w, struct {
		Sidebar template.HTML
	}{
		Sidebar: buildSidebar("Create"),
	})
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
			Sidebar template.HTML
			Header string
			Data   []struct {
				Name string
				Age  string
			}
		}{
			Sidebar: buildSidebar("Read"),
			Header: "Data List",
			Data:   data,
		})

		return
	}

	pageNotFoundHandler(w, r)
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
