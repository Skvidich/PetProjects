package storage

import "dataProcessor/common"

type Repository interface {
	StoreReport(overview *common.Report, components []common.ComponentMetrics) error
	StoreIncident(stat *common.ServiceIncident) error
}
