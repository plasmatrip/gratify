package auth

import (
	"github.com/plasmatrip/gratify/internal/api"
)

type Auth struct {
	deps api.Dependencies
}

func NewAuthService(desp api.Dependencies) *Auth {
	return &Auth{
		deps: desp,
	}
}
