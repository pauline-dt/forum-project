package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() {
	var err error

	DB, err = sql.Open("sqlite", "./forum.db")
	if err != nil {
		log.Fatal("Erreur ouverture base de données :", err)
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		image_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	createPostCategoriesTable := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY(post_id, category_id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);`

	queries := []string{
		createUsersTable,
		createPostsTable,
		createCategoriesTable,
		createPostCategoriesTable,
	}

	for _, query := range queries {
		_, err = DB.Exec(query)
		if err != nil {
			log.Fatal("Erreur création table :", err)
		}
	}

	insertDefaultCategories := `
	INSERT OR IGNORE INTO categories (name) VALUES
	('Général'),
	('Développement'),
	('Aide'),
	('Projet'),
	('Discussion');`

	_, err = DB.Exec(insertDefaultCategories)
	if err != nil {
		log.Fatal("Erreur insertion catégories :", err)
	}

	log.Println("Base de données initialisée")
}