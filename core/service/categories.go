package service

import (
	"expense-tracker-be/dto"
	"expense-tracker-be/storage"
	"fmt"

	"github.com/google/uuid"
)

func GetCategoriesAndSubCategories(userID uuid.UUID) []dto.CategoryWithSubcategories {
	fmt.Println(userID)
	categories, err := storage.GetCategoriesByUserID(userID)
	categoriesWithSubDtos := []dto.CategoryWithSubcategories{}
	if err != nil {
		fmt.Println("err fetching categories")
		return categoriesWithSubDtos
	}
	for _, category := range categories {
		subcategories, err := storage.GetSubcategoriesByCategoryID(category.ID)
		if err != nil {
			fmt.Println("err fetching subcategories")
			return categoriesWithSubDtos
		}
		subCategoryDtos := []dto.SubcategoryDetail{}
		for _, subCategory := range subcategories {
			subCategoryDto := dto.SubcategoryDetail{
				ID:   subCategory.ID,
				Name: subCategory.Name,
			}
			subCategoryDtos = append(subCategoryDtos, subCategoryDto)
		}
		categoryDto := dto.CategoryWithSubcategories{
			ID:            category.ID,
			Name:          category.Name,
			Subcategories: subCategoryDtos,
		}

		categoriesWithSubDtos = append(categoriesWithSubDtos, categoryDto)
	}

	return categoriesWithSubDtos
}
