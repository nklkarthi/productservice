package tests

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"productservice/db"
	"productservice/handlers"
	"productservice/middleware"
	"productservice/models"
	"productservice/utils"
	"strings"
	"testing"
)

func setupTestDB() {
	// Initialize an in-memory SQLite database for testing
	var err error
	db.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.DB.AutoMigrate(&models.Product{})
}

func TestGetProducts(t *testing.T) {
	setupTestDB()

	req, _ := http.NewRequest("GET", "/api/products", nil)
	req.Header.Set("Authorization", "Bearer "+getUserToken())
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestAddProduct(t *testing.T) {
	setupTestDB()

	var jsonStr = []byte(`{"name":"Product1","price":12.34}`)
	req, _ := http.NewRequest("POST", "/api/products", strings.NewReader(string(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+getAdminToken())
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestAddProductNonAdmin(t *testing.T) {
	setupTestDB()

	var jsonStr = []byte(`{"name":"Product1","price":12.34}`)
	req, _ := http.NewRequest("POST", "/api/products", strings.NewReader(string(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+getUserToken())
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestGetProduct(t *testing.T) {
	setupTestDB()

	// Add a product
	var jsonStr = []byte(`{"name":"Product1","price":12.34}`)
	addReq, _ := http.NewRequest("POST", "/api/products", strings.NewReader(string(jsonStr)))
	addReq.Header.Set("Authorization", "Bearer "+getAdminToken())
	addResponse := executeRequest(addReq)

	checkResponseCode(t, http.StatusOK, addResponse.Code)

	// Get the added product
	req, _ := http.NewRequest("GET", "/api/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+getUserToken())
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetProductNotFound(t *testing.T) {
	setupTestDB()

	req, _ := http.NewRequest("GET", "/api/products/999", nil)
	req.Header.Set("Authorization", "Bearer "+getUserToken())
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetProductsNoAuth(t *testing.T) {
	setupTestDB()

	req, _ := http.NewRequest("GET", "/api/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")
	router.HandleFunc("/api/products", handlers.AddProduct).Methods("POST")
	router.HandleFunc("/api/products/{id:[0-9]+}", handlers.GetProduct).Methods("GET")
	router.Use(middleware.JwtAuthentication) // Applying middleware
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getAdminToken() string {
	// Creating a mock admin user
	adminUser := models.User{
		Username: "admin",
		Role:     "admin",
	}
	token, _ := utils.GenerateJWT(adminUser.Username, adminUser.Role)
	return token
}

func getUserToken() string {
	// Creating a mock regular user
	user := models.User{
		Username: "user",
		Role:     "user",
	}
	token, _ := utils.GenerateJWT(user.Username, user.Role)
	return token
}
