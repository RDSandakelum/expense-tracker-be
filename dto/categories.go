package dto

import "github.com/google/uuid"

type SubcategoryDetail struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CategoryWithSubcategories struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	Subcategories []SubcategoryDetail `json:"subcategories"`
}
