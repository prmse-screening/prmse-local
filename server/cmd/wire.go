package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"server/internal/data"
	"server/internal/handlers"
	"server/internal/rpc"
	"server/internal/schedule"
	"server/internal/services"
)

func wireApp() (*gin.Engine, error) {
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
