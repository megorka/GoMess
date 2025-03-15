package models

type Post struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"user_id" db:"user_id"`
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
	Created_at string `json:"created_at" db:"created_at"`
}
