package storage

import "dataCollector/pkg/types"

type Storage interface {
	StoreRawReport(stat *types.ServiceStatus) error
}
