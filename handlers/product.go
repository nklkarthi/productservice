package handlers

import (
	"encoding/json"
	"net/http"
	"productservice/models"
	"productservice/services"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") == nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	products := services.GetAllProducts()
	if len(products) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "No products found"})
		return
	}

	json.NewEncoder(w).Encode(products)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role != "admin" {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)

	err := services.AddProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product is saved successfully"})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") == nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := services.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}
