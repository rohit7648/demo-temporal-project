//go:build wireinject
// +build wireinject

package main

import (
	"demo-temporal-project/client/temporal"
	"demo-temporal-project/configs"
	"demo-temporal-project/internal/service"
	"demo-temporal-project/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*configs.Server, string, *configs.DemoTemporalProject, *configs.Temporal, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, temporal.ProviderSet, newApp))
}
