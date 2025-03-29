package notificator

import (
	"dataProcessor/common"
	"time"
)

type Notificator interface {
	StartIncident(name string, incdTime time.Time, comp common.Component) error
	EndIncident(incd *common.Incident) error
}
