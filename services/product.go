package services

import (
	"errors"
	"productservice/models"
)

var products []models.Product
var idCounter = 1

func GetAllProducts() []models.Product {
	return products
}

func AddProduct(product models.Product) error {
	// Check for uniqueness of product name
	for _, p := range products {
		if p.Name == product.Name {
			return errors.New("the Product name must be unique")
		}
	}

	product.ID = idCounter
	idCounter++
	products = append(products, product)
	return nil
}

func GetProductByID(id int) (models.Product, error) {
	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}
	return models.Product{}, errors.New("Product not found")
}
