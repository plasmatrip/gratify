package info

import "github.com/plasmatrip/gratify/internal/deps"

type Info struct {
	deps deps.Dependencies
}

func NewInfoService(deps deps.Dependencies) *Info {
	return &Info{
		deps: deps,
	}
}
