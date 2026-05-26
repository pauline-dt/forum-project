package main

import (
	"html/template"
	"log"
	"net/http"

	"forum-project/internal/database"

	"golang.org/x/crypto/bcrypt"
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
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html")
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email et mot de passe obligatoires", http.StatusBadRequest)
			return
		}

		var id int
		var username string
		var hashedPassword string

		err := database.DB.QueryRow(
			"SELECT id, username, password FROM users WHERE email = ?",
			email,
		).Scan(&id, &username, &hashedPassword)

		if err != nil {
			http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "register.html")
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur hash mot de passe", http.StatusInternalServerError)
			return
		}

		_, err = database.DB.Exec(
			"INSERT INTO users(username, email, password) VALUES(?, ?, ?)",
			username,
			email,
			string(hashedPassword),
		)

		if err != nil {
			http.Error(w, "Email déjà utilisé", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
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