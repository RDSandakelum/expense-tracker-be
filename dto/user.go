package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID           uint
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
}
