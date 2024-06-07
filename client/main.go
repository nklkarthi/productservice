package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "http://localhost:8000/api"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	// Login as non-admin user
	userToken, err := login("user", "userpass")
	if err != nil {
		log.Fatalf("Login failed for user: %v", err)
	}

	// Try to add a product as non-admin user
	err = addProduct(userToken, Product{Name: "NonAdminProduct", Price: 123.45})
	if err == nil {
		log.Fatalf("Expected error when adding product as non-admin, but got none")
	} else {
		fmt.Println("Non-admin user add product error:", err)
	}

	// Login as admin user
	adminToken, err := login("admin", "adminpass")
	if err != nil {
		log.Fatalf("Login failed for admin: %v", err)
	}

	// Add a product as admin user
	product := Product{Name: "AdminProduct", Price: 123.45}
	err = addProduct(adminToken, product)
	if err != nil {
		log.Fatalf("Failed to add product as admin: %v", err)
	}
	fmt.Println("Product added successfully as admin")

	// Attempt to add the same product again to check for uniqueness
	err = addProduct(adminToken, product)
	if err != nil {
		fmt.Println("Duplicate product error:", err)
	}

	// Get all products as a privileged user
	products, err := getProducts(userToken)
	if err != nil {
		log.Fatalf("Failed to get products: %v", err)
	}
	if len(products) == 0 {
		fmt.Println("No products found")
	} else {
		fmt.Println("Products:", products)
	}

	// Get product by ID as a privileged user
	productID := 1
	productDetails, err := getProduct(userToken, productID)
	if err != nil {
		fmt.Printf("Failed to get product ID %d: %v\n", productID, err)
	} else {
		fmt.Printf("Product ID %d details: %+v\n", productID, productDetails)
	}

	// Anonymous user trying to get products
	err = getProductsNoAuth()
	if err != nil {
		fmt.Println("Anonymous user get products error:", err)
	}
}

func login(username, password string) (string, error) {
	credentials := Credentials{Username: username, Password: password}
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status: %s", resp.Status)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["token"], nil
}

func addProduct(token string, product Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/products", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to add product: %s", bodyString)
	}

	return nil
}

func getProducts(token string) ([]Product, error) {
	req, err := http.NewRequest("GET", baseURL+"/products", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get products: %s", resp.Status)
	}

	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func getProduct(token string, id int) (Product, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products/%d", baseURL, id), nil)
	if err != nil {
		return Product{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Product{}, fmt.Errorf("failed to get product: %s", resp.Status)
	}

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func getProductsNoAuth() error {
	req, err := http.NewRequest("GET", baseURL+"/products", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get products: %s", resp.Status)
	}

	return nil
}
