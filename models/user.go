package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
