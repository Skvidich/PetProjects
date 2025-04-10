package server

import (
	"dataCollector/internal/core/coordinator"
	"dataCollector/internal/core/storage"
	"dataCollector/internal/grpcAPI"
	"dataCollector/internal/logger"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	gRPCServer *grpc.Server
	errLog     logger.Logger
	addr       string
}

func New(
	cord *coordinator.Coordinator,
	errLog logger.Logger,
	stor storage.Storage,
	addr string,
) *App {
	gRPCServer := grpc.NewServer()
	cntrl := grpcAPI.NewControlAPI(cord)
	storApi := grpcAPI.NewStorageAPI(stor)
	grpcAPI.Register(gRPCServer, cntrl, storApi)

	return &App{
		gRPCServer: gRPCServer,
		errLog:     errLog,
		addr:       addr,
	}
}

func (a *App) Run() error {

	l, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	if err = a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
