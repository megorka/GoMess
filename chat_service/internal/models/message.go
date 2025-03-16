package models

type Message struct {
	ID      int    `json:"id"`
	ToID    int    `json:"to_id"`
	Content string `json:"content"`
}
