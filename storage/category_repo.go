package storage

import "github.com/google/uuid"

func GetCategoriesByUserID(userID uuid.UUID) ([]Category, error) {
	var categories []Category
	result := DB.Where("user_id = ?", userID).Find(&categories)
	return categories, result.Error
}

func GetCategoryByID(id uuid.UUID) (*Category, error) {
	var category Category
	result := DB.Where("id = ?", id).First(&category)
	return &category, result.Error
}

func GetSubcategoriesByCategoryID(categoryID uuid.UUID) ([]SubCategory, error) {
	var subcategories []SubCategory
	result := DB.Where("category_id = ?", categoryID).Find(&subcategories)
	return subcategories, result.Error
}

func GetSubcategoryByID(id uuid.UUID) (*SubCategory, error) {
	var subCategory SubCategory
	result := DB.First(&subCategory, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &subCategory, nil
}

func CreateSubcategory(categoryID uuid.UUID, name string) (*SubCategory, error) {
	subCategory := SubCategory{
		ID:         uuid.New(),
		CategoryID: categoryID,
		Name:       name,
	}
	result := DB.Create(&subCategory)
	return &subCategory, result.Error
}
