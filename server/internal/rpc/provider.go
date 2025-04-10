package rpc

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server/internal/config"
)

var ProviderSet = wire.NewSet(NewRpcClient)

func NewRpcClient() *[]WorkerClient {
	clients := make([]WorkerClient, 0, len(config.Cfg.Worker.Endpoints))
	for _, endpoint := range config.Cfg.Worker.Endpoints {
		conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("failed to connect to gRPC server: %v", err)
		}
		c := NewWorkerClient(conn)
		clients = append(clients, c)
	}
	return &clients
}
