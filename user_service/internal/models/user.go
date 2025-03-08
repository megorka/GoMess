package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	LastName string `json:"lastname" db:"lastname"`
	Email    string `json:"email" db:"email"`
}
