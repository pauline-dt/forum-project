package main

import (
	"html/template"
	"log"
	"net/http"

	"forum-project/internal/database"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl, nil)

	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "index.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html")
}

func main() {

	database.InitDatabase()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	log.Println("Serveur lancé sur http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}