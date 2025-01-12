package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/plasmatrip/gratify/internal/api"
	"github.com/plasmatrip/gratify/internal/models"
)

type Auth struct {
	deps api.Dependencies
}

func NewAuthService(desp api.Dependencies) *Auth {
	return &Auth{
		deps: desp,
	}
}

func (a *Auth) MakeLoginToken(lr models.LoginRequest) (string, error) {
	payload := jwt.MapClaims{
		"sub": lr.Login,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(a.deps.Config.TokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
