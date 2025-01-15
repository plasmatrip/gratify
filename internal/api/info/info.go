package info

import "github.com/plasmatrip/gratify/internal/api"

type Info struct {
	deps api.Dependencies
}

func NewInfoService(deps api.Dependencies) *Info {
	return &Info{
		deps: deps,
	}
}
