package grpcAPI

import (
	"dataCollector/pkg/types"
	collectorV1 "github.com/Skvidich/CollectorProto/gen/go"
	"google.golang.org/grpc"
)

func Register(gRPC *grpc.Server, control *ControlAPI, storage *StorageAPI) {
	collectorV1.RegisterControlServer(gRPC, control)
	collectorV1.RegisterStorageServer(gRPC, storage)
}

type Coordinator interface {
	Getter(name string) (types.GetterInfo, error)

	GetterList() []types.GetterInfo

	Start(name string) error

	StartAll() error

	StopAll()

	Stop(name string) error

	Shutdown()
}

type Storage interface {
}
