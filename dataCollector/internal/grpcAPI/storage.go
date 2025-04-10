package grpcAPI

import collectorV1 "github.com/Skvidich/CollectorProto/gen/go"

type StorageAPI struct {
	collectorV1.UnimplementedStorageServer
	stor Storage
}

func NewStorageAPI(stor Storage) *StorageAPI {
	return &StorageAPI{
		UnimplementedStorageServer: collectorV1.UnimplementedStorageServer{},
		stor:                       stor,
	}
}
