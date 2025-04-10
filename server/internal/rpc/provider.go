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
	clients := make([]WorkerClient, 0, len(config.Cfg.Workers.Endpoints))
	for _, endpoint := range config.Cfg.Workers.Endpoints {
		conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect to gRPC server: %v", err)
		}
		c := NewWorkerClient(conn)
		clients = append(clients, c)
	}
	return &clients
}
