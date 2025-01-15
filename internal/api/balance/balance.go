package balance

import (
	"github.com/plasmatrip/gratify/internal/api"
)

type Balance struct {
	deps api.Dependencies
}

func NewBalanceService(deps api.Dependencies) *Balance {
	return &Balance{
		deps: deps,
	}
}
