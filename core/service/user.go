package service

import (
	"expense-tracker-be/core/domain"
	"expense-tracker-be/storage"
)

func GetUserByEmail(email string) *domain.User {
	user, err := storage.GetUserByEmail(email)
	if err != nil {
		return nil
	}
	return &domain.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: user.Password,
		Email:        user.Email,
	}
}
