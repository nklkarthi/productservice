package tests

import (
	"net/http"
	"net/http/httptest"
	"productservice/handlers"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestLogin(t *testing.T) {
	var jsonStr = []byte(`{"username":"admin","password":"adminpass"}`)
	req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(string(jsonStr)))
	response := executeAuthRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeAuthRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.ServeHTTP(rr, req)

	return rr
}
