package graph

import (
	"GRAPHQL/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services service.Service
}

func NewResolver(services service.Service) *Resolver {
	return &Resolver{services: services}
}
