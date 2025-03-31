package processor

import (
	"context"
	"dataProcessor/internal/core/alert"
	"dataProcessor/internal/core/reader"
	"dataProcessor/internal/core/storage"
	"dataProcessor/pkg/errors"
	"dataProcessor/pkg/models"
	"time"
)

const AggregationInterval = time.Second * 60

type incidentKey struct {
	service   string
	component string
}
type incidentValue struct {
	status    string
	startTime time.Time
}

type ServiceMetrics struct {
	service       string
	intervalStart time.Time
	intervalEnd   time.Time
	components    map[string]map[string]int // component -> status -> count
}

type MetricsBuffer map[string]*ServiceMetrics
type IncidentBuffer map[incidentKey]incidentValue
type StatusEngine struct {
	activeIncidents IncidentBuffer
	currentMetrics  MetricsBuffer
	storage         storage.Repository
	statusQueue     reader.Reader
	alertService    alert.Handler
	ctx             context.Context
	cancel          context.CancelFunc
	ErrChan         chan errors.TaggedError
}

func NewEngine(
	saver storage.Repository,
	consumer reader.Reader,
	notificator alert.Handler,
) (*StatusEngine, error) {

	ctx, cancel := context.WithCancel(context.Background())

	return &StatusEngine{
		activeIncidents: make(map[incidentKey]incidentValue),
		currentMetrics:  make(map[string]*ServiceMetrics),
		storage:         saver,
		statusQueue:     consumer,
		alertService:    notificator,
		ctx:             ctx,
		cancel:          cancel,
		ErrChan:         make(chan errors.TaggedError),
	}, nil
}

func (proc *StatusEngine) Start(timeout time.Duration) {

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	for {
		select {
		case <-proc.ctx.Done():
			return
		case <-ctx.Done():
			mess, err := proc.statusQueue.Next()
			if err != nil {
				proc.sendError(errors.TaggedError{Type: errors.TypeConsumer, Err: err})
				return
			}

			if mess != nil {

				errInf := proc.process(mess)
				if errInf.Err != nil {
					proc.sendError(errInf)
					return
				}
			} else {
				ctx, _ = context.WithTimeout(context.Background(), timeout)
			}

		}
	}

}

func (proc *StatusEngine) process(update *models.ServiceStatus) errors.TaggedError {

	if err := proc.trackIncidents(update); err.Err != nil {
		return err
	}

	if err := proc.aggregateMetrics(update); err.Err != nil {
		return err
	}

	return errors.TaggedError{Err: nil}
}

func (proc *StatusEngine) trackIncidents(status *models.ServiceStatus) errors.TaggedError {
	var err error

	for _, component := range status.Components {
		key := incidentKey{
			service:   status.Name,
			component: component.Name,
		}
		incident, ok := proc.activeIncidents[key]

		if ok && models.IsUp(component.Status) {
			endedIncident := &models.ServiceIncident{
				Name:      key.service,
				Component: key.component,
				State:     incident.status,
				StartTime: incident.startTime,
				EndTime:   status.Time,
			}
			err = proc.storage.StoreIncident(endedIncident)
			if err != nil {
				return errors.TaggedError{Type: errors.TypeSaver, Err: err}
			}

			err = proc.alertService.NotifyEnd(endedIncident)
			if err != nil {
				return errors.TaggedError{Type: errors.TypeNotificator, Err: err}
			}
			delete(proc.activeIncidents, key)
		}

		if !ok && models.IsDown(component.Status) {
			proc.activeIncidents[incidentKey{service: status.Name, component: component.Name}] = incidentValue{startTime: status.Time, status: component.Status}
			err = proc.alertService.NotifyStart(status.Name, status.Time, component)
			if err != nil {
				return errors.TaggedError{Type: errors.TypeNotificator, Err: err}
			}
		}
	}
	return errors.TaggedError{Err: nil}
}

func (proc *StatusEngine) aggregateMetrics(mess *models.ServiceStatus) errors.TaggedError {

	_, ok := proc.currentMetrics[mess.Name]
	if !ok {
		proc.currentMetrics[mess.Name] = &ServiceMetrics{service: mess.Name, intervalStart: mess.Time, components: make(map[string]map[string]int)}
	}

	for _, component := range mess.Components {
		if _, exists := proc.currentMetrics[mess.Name].components[component.Name]; !exists {
			proc.currentMetrics[mess.Name].components[component.Name] = make(map[string]int)
		}
		proc.currentMetrics[mess.Name].components[component.Name][component.Status]++
	}

	deadline := proc.currentMetrics[mess.Name].intervalStart.Add(AggregationInterval)
	if mess.Time.After(deadline) {
		proc.currentMetrics[mess.Name].intervalEnd = mess.Time
		overview, components := parseMetrics(proc.currentMetrics[mess.Name])
		err := proc.storage.StoreReport(overview, components)
		if err != nil {
			return errors.TaggedError{Type: errors.TypeSaver, Err: err}
		}
		delete(proc.currentMetrics, mess.Name)
	}

	return errors.TaggedError{Err: nil}
}

func (proc *StatusEngine) sendError(errInf errors.TaggedError) {
	proc.ErrChan <- errInf
}

func (proc *StatusEngine) GetBackup() (IncidentBuffer, MetricsBuffer) {
	return proc.activeIncidents, proc.currentMetrics
}

func (proc *StatusEngine) Close() {
	proc.cancel()
	close(proc.ErrChan)
}

func parseMetrics(metrics *ServiceMetrics) (*models.Report, []models.ComponentMetrics) {
	report := &models.Report{
		Name:        metrics.service,
		PeriodStart: metrics.intervalStart,
		PeriodEnd:   metrics.intervalEnd,
	}

	componentMetrics := make([]models.ComponentMetrics, 0)

	for name, stat := range metrics.components {
		for status, count := range stat {
			componentMetrics = append(componentMetrics, models.ComponentMetrics{
				Component: name,
				State:     status,
				Count:     count,
			})
		}
	}
	return report, componentMetrics
}
