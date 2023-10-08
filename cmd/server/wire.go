//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"zerocmf/configs"
	"zerocmf/internal/server"

	"github.com/google/wire"
)

// wireApp init application.
func wireApp(configs.Config) *ServiceContext {
	panic(wire.Build(server.ProviderSet, newApp))
}
