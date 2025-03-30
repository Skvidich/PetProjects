package alert

import (
	"dataProcessor/pkg/models"
	"fmt"
	"time"
)

type IncidentReporter struct {
}

func NewIncidentNotificator() *IncidentReporter {
	return &IncidentReporter{}
}

func (r *IncidentReporter) NotifyStart(name string, incdTime time.Time, comp models.Component) error {
	fmt.Println("incident started")
	return nil
}

func (r *IncidentReporter) NotifyEnd(incident *models.ServiceIncident) error {
	fmt.Println("incident ended")
	return nil
}
