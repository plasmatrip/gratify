package orders

import (
	"github.com/plasmatrip/gratify/internal/controller"
	"github.com/plasmatrip/gratify/internal/deps"
)

type Orders struct {
	deps       deps.Dependencies
	controller *controller.Controller
}

func NewOrdersService(deps deps.Dependencies, controller *controller.Controller) *Orders {
	return &Orders{
		deps:       deps,
		controller: controller,
	}
}
