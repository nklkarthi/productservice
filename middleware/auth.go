package middleware

import (
	"context"
	"net/http"
	"productservice/utils"
	"strings"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		tokenParts := strings.Split(tokenHeader, " ")
		if len(tokenParts) != 2 {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		token := tokenParts[1]
		claims, err := utils.DecodeJWT(token)
		if err != nil {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
