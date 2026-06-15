package models

type Post struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Author   string    `json:"author"`
	Category string    `json:"category"`
	Image    string    `json:"image"`
	Likes    int       `json:"likes"`
	Dislikes int       `json:"dislikes"`
	Comments []Comment `json:"comments"`
}