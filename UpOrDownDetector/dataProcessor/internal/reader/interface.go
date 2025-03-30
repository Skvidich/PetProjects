package reader

import "dataProcessor/pkg/models"

type Reader interface {
	Next() (*models.ServiceStatus, error)
}
