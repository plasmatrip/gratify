package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (a Auth) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.deps.Logger.Sugar.Info("missing authorization header")
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.deps.Logger.Sugar.Infow("invalid authorization header format", "parts", parts)
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(a.deps.Config.TokenSecret), nil
			})

		if err != nil {
			a.deps.Logger.Sugar.Infow("JWT token error", "error", err)
			http.Error(w, "JWT token error", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			a.deps.Logger.Sugar.Infow("invalid token", "token", token)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ValidLogin{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
