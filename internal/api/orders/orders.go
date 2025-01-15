package orders

import (
	"github.com/plasmatrip/gratify/internal/api"
)

type Orders struct {
	deps api.Dependencies
}

func NewOrdersService(deps api.Dependencies) *Orders {
	return &Orders{
		deps: deps,
	}
}
