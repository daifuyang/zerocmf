//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"zerocmf/configs"
	"zerocmf/internal/biz"
	"zerocmf/internal/data"
	"zerocmf/internal/server"
	"zerocmf/internal/service"

	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configs.Config) (App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
