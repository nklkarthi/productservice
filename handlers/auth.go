package handlers

import (
	"encoding/json"
	"net/http"
	"productservice/models"
	"productservice/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	token, err := services.Authenticate(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
