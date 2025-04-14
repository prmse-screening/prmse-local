//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
	"server/internal/data"
	"server/internal/handlers"
	"server/internal/rpc"
	"server/internal/schedule"
	"server/internal/services"
)

func wireApp() (*http.Server, error) {
	wire.Build(
		data.ProviderSet,
		rpc.ProviderSet,
		schedule.ProviderSet,
		services.ServiceSet,
		handlers.ProviderSet,
		NewServer,
	)
	return nil, nil
}
