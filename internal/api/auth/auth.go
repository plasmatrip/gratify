package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/plasmatrip/gratify/internal/api"
	"github.com/plasmatrip/gratify/internal/models"
)

type ValidLogin struct {
}

type Claims struct {
	jwt.StandardClaims
	UserdID int32
}

type Auth struct {
	deps api.Dependencies
}

func NewAuthService(desp api.Dependencies) *Auth {
	return &Auth{
		deps: desp,
	}
}

func (a *Auth) LoginToken(lr models.LoginRequest) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Subject:   lr.Login,
		},
		UserdID: lr.ID,
	}
	// claims := jwt.MapClaims{
	// 	"sub": lr.Login,
	// 	"exp": time.Now().Add(time.Hour * 72).Unix(),
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(a.deps.Config.TokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
