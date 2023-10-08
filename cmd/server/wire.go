//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"zerocmf/internal/data"
	"zerocmf/internal/server"
	"zerocmf/internal/svc"

	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*svc.ServiceContext) (App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, newApp))
}
