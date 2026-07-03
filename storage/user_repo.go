package storage

import "github.com/google/uuid"

func GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
