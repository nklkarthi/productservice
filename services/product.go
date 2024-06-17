package services

import (
	"errors"
	"gorm.io/gorm"
	"productservice/db"
	"productservice/models"
)

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := db.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func AddProduct(product models.Product) error {
	result := db.DB.Create(&product)
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: products.name" {
			return errors.New("the Product name must be unique")
		}
		return result.Error
	}
	return nil
}

func GetProductByID(id uint) (models.Product, error) {
	var product models.Product
	result := db.DB.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return product, errors.New("Product not found")
		}
		return product, result.Error
	}
	return product, nil
}
