package orders

import (
	"github.com/plasmatrip/gratify/internal/deps"
)

type Orders struct {
	deps deps.Dependencies
}

func NewOrdersService(deps deps.Dependencies) *Orders {
	return &Orders{
		deps: deps,
	}
}
