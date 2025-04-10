package storage

import "dataProcessor/pkg/models"

type Repository interface {
	StoreReport(overview *models.Report, components []models.ComponentMetrics) error
	StoreIncident(stat *models.ServiceIncident) error
}
