package balance

import "github.com/plasmatrip/gratify/internal/deps"

type Balance struct {
	deps deps.Dependencies
}

func NewBalanceService(deps deps.Dependencies) *Balance {
	return &Balance{
		deps: deps,
	}
}
