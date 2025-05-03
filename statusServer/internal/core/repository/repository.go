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
type ReportRepository interface {
	// GetReportByID возвращает отчёт по его ID.
	GetReportByID(ctx context.Context, id int) (*models.Report, Error)
	// ListReports возвращает все отчёты (с пагинацией).
	ListReports(ctx context.Context, offset, limit int) ([]*models.Report, Error)
	// ListReportsByService возвращает отчёты для заданного сервиса.
	ListReportsByService(ctx context.Context, serviceName string, offset, limit int) ([]*models.Report, Error)
	// ListReportsByPeriod возвращает отчёты за указанный диапазон времени.
	ListReportsByPeriod(ctx context.Context, from, to time.Time, offset, limit int) ([]*models.Report, Error)
	// CountReports возвращает общее число отчётов (для пагинации).
	CountReports(ctx context.Context) (int, Error)
	// CountReportsByService возвращает число отчётов для конкретного сервиса.
	CountReportsByService(ctx context.Context, serviceName string) (int, Error)
}

type ComponentMetricRepository interface {
	// GetMetricByID возвращает метрику по её ID.
	GetMetricByID(ctx context.Context, id int) (*models.ComponentMetric, Error)
	// ListMetricsByReport возвращает все метрики, относящиеся к данному отчёту.
	ListMetricsByReport(ctx context.Context, reportID int) ([]*models.ComponentMetric, Error)
	// ListMetricsByComponent возвращает метрики для указанного компонента.
	ListMetricsByComponent(ctx context.Context, comp string) ([]*models.ComponentMetric, Error)

	// ListMetricsByStatus возвращает метрики с данным статусом.
	ListMetricsByStatus(ctx context.Context, status string) ([]*models.ComponentMetric, Error)
	// ListMetricsByReportAndStatus – комбинированный фильтр по отчёту и статусу.
	ListMetricsByReportAndStatus(ctx context.Context, reportID int, status string) ([]*models.ComponentMetric, Error)
	// CountMetricsByReport возвращает количество метрик в отчёте.
	CountMetricsByReport(ctx context.Context, reportID int) (int, Error)
}

type IncidentRepository interface {
	// GetIncidentByID возвращает инцидент по его ID.
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, Error)
	// ListIncidents возвращает список инцидентов за период, с фильтрами по сервису и компоненту.
	ListIncidents(ctx context.Context, service, component string, from, to time.Time) ([]*models.Incident, Error)

	// ListIncidentsByService возвращает все инциденты для конкретного сервиса.
	ListIncidentsByService(ctx context.Context, service string, from, to time.Time) ([]*models.Incident, Error)
	// ListIncidentsByComponent возвращает все инциденты для конкретного компонента.
	ListIncidentsByComponent(ctx context.Context, component string, from, to time.Time) ([]*models.Incident, Error)
	// CountIncidentsByService возвращает число инцидентов для сервиса за период.
	CountIncidentsByService(ctx context.Context, service string, from, to time.Time) (int, Error)
	// CountIncidentsByComponent возвращает число инцидентов для компонента за период.
	CountIncidentsByComponent(ctx context.Context, component string, from, to time.Time) (int, Error)
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

/*
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
*/
