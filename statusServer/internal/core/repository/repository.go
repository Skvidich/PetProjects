package repository

import (
	"context"
	"fmt"
	"statusServer/pkg/models"
	"time"
)

type StatusRepository interface {
	ReportRepository
	IncidentRepository
	ComponentMetricRepository
}
type Error struct {
	Err     error
	Problem string
	Method  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Method %s encountered error: %s %s", e.Method, e.Problem, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}

type ReportRepository interface {
	GetReportByID(ctx context.Context, id int) (*models.Report, error)
	ListReports(ctx context.Context, offset, limit int) ([]*models.Report, error)
	ListReportsByService(ctx context.Context, serviceName string, offset, limit int) ([]*models.Report, error)
	ListReportsByPeriod(ctx context.Context, from, to time.Time, offset, limit int) ([]*models.Report, error)
	CountReports(ctx context.Context) (int, error)
	CountReportsByService(ctx context.Context, serviceName string) (int, error)
}

type ComponentMetricRepository interface {
	GetMetricByID(ctx context.Context, id int) (*models.ComponentMetric, error)
	ListMetricsByReport(ctx context.Context, reportID int) ([]*models.ComponentMetric, error)
	ListMetricsByComponent(ctx context.Context, comp string) ([]*models.ComponentMetric, error)
	ListMetricsByStatus(ctx context.Context, status string) ([]*models.ComponentMetric, error)
	ListMetricsByReportAndStatus(ctx context.Context, reportID int, status string) ([]*models.ComponentMetric, error)
	CountMetricsByReport(ctx context.Context, reportID int) (int, error)
}

type IncidentRepository interface {
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, error)
	ListIncidents(ctx context.Context, service, component string, from, to time.Time) ([]*models.Incident, error)
	ListIncidentsByService(ctx context.Context, service string, from, to time.Time) ([]*models.Incident, error)
	ListIncidentsByComponent(ctx context.Context, component string, from, to time.Time) ([]*models.Incident, error)
	CountIncidentsByService(ctx context.Context, service string, from, to time.Time) (int, error)
	CountIncidentsByComponent(ctx context.Context, component string, from, to time.Time) (int, error)
}
