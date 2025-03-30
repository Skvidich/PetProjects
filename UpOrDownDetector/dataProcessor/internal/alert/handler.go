package alert

import (
	"dataProcessor/pkg/models"
	"time"
)

type Handler interface {
	NotifyStart(name string, incdTime time.Time, comp models.Component) error
	NotifyEnd(incident *models.ServiceIncident) error
}
