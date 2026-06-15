package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum-project/internal/database"
	"forum-project/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
type PageData struct {
	IsLoggedIn bool
	Posts      []models.Post
}
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

	_, isLoggedIn := getUserIDFromCookie(r)

	rows, err := database.DB.Query(
		"SELECT id, title, content FROM posts ORDER BY id DESC",
	)

	if err != nil {
		http.Error(w, "Erreur récupération posts", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
		)

		if err != nil {
			continue
		}

		posts = append(posts, post)
	}

	data := PageData{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)

	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
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

		cookie := &http.Cookie{
			Name:     "session_user_id",
			Value:    fmt.Sprintf("%d", id),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)

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

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	cookie := &http.Cookie{
		Name:     "session_user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUserIDFromCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_user_id")
	if err != nil {
		return "", false
	}

	if cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, isLoggedIn := getUserIDFromCookie(r)

	if !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		renderTemplate(w, "create_post.html")
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID := r.FormValue("category")

		if title == "" || content == "" || categoryID == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}

		result, err := database.DB.Exec(
			"INSERT INTO posts(user_id, title, content) VALUES(?, ?, ?)",
			userID,
			title,
			content,
		)

		if err != nil {
			http.Error(w, "Erreur création du post", http.StatusInternalServerError)
			return
		}

		postID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Erreur récupération du post", http.StatusInternalServerError)
			return
		}

		_, err = database.DB.Exec(
			"INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)",
			postID,
			categoryID,
		)

		if err != nil {
			http.Error(w, "Erreur association catégorie", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}
func apiPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT 
			posts.id,
			posts.title,
			posts.content,
			users.username,
			IFNULL(categories.name, ''),
			IFNULL(posts.image_path, ''),
			0,
			0
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN post_categories ON posts.id = post_categories.post_id
		LEFT JOIN categories ON post_categories.category_id = categories.id
		ORDER BY posts.created_at DESC
	`)

	if err != nil {
		http.Error(w, "Erreur récupération posts", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&post.Category,
			&post.Image,
			&post.Likes,
			&post.Dislikes,
		)

		if err != nil {
			http.Error(w, "Erreur lecture post", http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
func apiAddCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, isLoggedIn := getUserIDFromCookie(r)
	if !isLoggedIn {
		http.Error(w, "Vous devez être connecté", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if postID == "" || content == "" {
		http.Error(w, "Commentaire vide", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(
		"INSERT INTO comments(post_id, user_id, content) VALUES(?, ?, ?)",
		postID,
		userID,
		content,
	)

	if err != nil {
		http.Error(w, "Erreur ajout commentaire", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {

	database.InitDatabase()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/create-post", createPostHandler)
	http.HandleFunc("/api/posts", apiPostsHandler)
	http.HandleFunc("/api/comments/add", apiAddCommentHandler)

	log.Println("Serveur lancé sur http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}