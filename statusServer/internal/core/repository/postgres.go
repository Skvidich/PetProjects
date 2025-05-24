package repository

import (
	"context"
	"database/sql"
	"statusServer/pkg/models"
	"time"
)

type PostgresRepo struct {
	// TODO add comprehensive logic
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (p *PostgresRepo) Close() error {
	return nil
}

// Component interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetMetricByID(ctx context.Context, id int) (*models.ComponentMetric, error) {

	op := "GetMetricsByID"
	query := "SELECT id,report_id,component_name,status,status_count FROM component_metrics WHERE id=$1"

	res := p.db.QueryRowContext(ctx, query, id)
	if res.Err() != nil {
		return nil, &Error{
			Err:     res.Err(),
			Problem: "can't execute query",
			Method:  op,
		}
	}

	var metric models.ComponentMetric

	err := res.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)

	if err != nil {
		return nil, &Error{
			Err:     err,
			Problem: "can't scan row",
			Method:  op,
		}
	}

	return &metric, nil

}

func (p *PostgresRepo) ListMetricsByReport(ctx context.Context, reportID int) ([]*models.ComponentMetric, error) {

	op := "ListMetricsByReport"
	query := "SELECT id,report_id,component_name,status,status_count FROM component_metrics  WHERE report_id=$1"

	rows, err := p.db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, &Error{
			Err:     err,
			Problem: "can't execute query",
			Method:  op,
		}
	}
	defer rows.Close()
	metrics := make([]*models.ComponentMetric, 0)

	for rows.Next() {
		var metric models.ComponentMetric
		err = rows.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

func (p *PostgresRepo) ListMetricsByComponent(ctx context.Context, comp string) ([]*models.ComponentMetric, error) {

	op := "ListMetricsByComponent"
	query := "SELECT id,report_id,component_name,status,status_count FROM component_metrics  WHERE component_name=$1"

	rows, err := p.db.QueryContext(ctx, query, comp)
	if err != nil {
		return nil, &Error{
			Err:     err,
			Problem: "can't execute query",
			Method:  op,
		}
	}
	defer rows.Close()
	metrics := make([]*models.ComponentMetric, 0)

	for rows.Next() {
		var metric models.ComponentMetric
		err = rows.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil

}

func (p *PostgresRepo) ListMetricsByStatus(ctx context.Context, status string) ([]*models.ComponentMetric, error) {
	op := "ListMetricsByStatus"
	query := "SELECT id,report_id,component_name,status,status_count FROM component_metrics  WHERE status=$1"

	rows, err := p.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, &Error{
			Err:     err,
			Problem: "can't execute query",
			Method:  op,
		}
	}
	defer rows.Close()
	metrics := make([]*models.ComponentMetric, 0)

	for rows.Next() {
		var metric models.ComponentMetric
		err = rows.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil

}

func (p *PostgresRepo) ListMetricsByReportAndStatus(ctx context.Context, reportID int, status string) ([]*models.ComponentMetric, error) {
	op := "ListMetricsByReportAndStatus"
	query := "SELECT id,report_id,component_name,status,status_count FROM component_metrics  WHERE status=$1 AND report_id=$2"

	rows, err := p.db.QueryContext(ctx, query, status, reportID)
	if err != nil {
		return nil, &Error{
			Err:     err,
			Problem: "can't execute query",
			Method:  op,
		}
	}
	defer rows.Close()
	metrics := make([]*models.ComponentMetric, 0)

	for rows.Next() {
		var metric models.ComponentMetric
		err = rows.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

func (p *PostgresRepo) CountMetricsByReport(ctx context.Context, reportID int) (int, error) {
	op := "CountMetricsByReport"
	query := "SELECT COUNT(*) FROM component_metrics  WHERE report_id=$1"

	res := p.db.QueryRowContext(ctx, query, reportID)
	if res.Err() != nil {
		return 0, &Error{
			Err:     res.Err(),
			Problem: "can't execute query",
			Method:  op,
		}
	}

	var metricCount int

	err := res.Scan(&metricCount)

	if err != nil {
		return 0, &Error{
			Err:     err,
			Problem: "can't scan row",
			Method:  op,
		}
	}

	return metricCount, nil
}

// Incident interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetIncidentByID(ctx context.Context, id int) (*models.Incident, error) {
	op := "GetIncidentByID"
	query := "SELECT id, service_name, component_name, status, start_time, end_time FROM incidents WHERE id=$1"

	row := p.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var inc models.Incident
	if err := row.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
		return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return &inc, nil
}

func (p *PostgresRepo) ListIncidents(ctx context.Context, service, component string, from, to time.Time) ([]*models.Incident, error) {
	op := "ListIncidents"
	query := " SELECT id, service_name, component_name, status, start_time, end_time FROM incidents WHERE service_name=$1 AND component_name=$2 AND start_time>=$3 AND end_time<=$4 ORDER BY id"
	rows, err := p.db.QueryContext(ctx, query, service, component, from, to)
	if err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, nil
}

func (p *PostgresRepo) ListIncidentsByService(ctx context.Context, service string, from, to time.Time) ([]*models.Incident, error) {
	op := "ListIncidentsByService"
	query := "SELECT id, service_name, component_name, status, start_time, end_time FROM incidents WHERE service_name=$1 AND start_time>=$2 AND end_time<=$3 ORDER BY id"
	rows, err := p.db.QueryContext(ctx, query, service, from, to)
	if err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, nil
}

func (p *PostgresRepo) ListIncidentsByComponent(ctx context.Context, component string, from, to time.Time) ([]*models.Incident, error) {
	op := "ListIncidentsByComponent"
	query := "SELECT id, service_name, component_name, status, start_time, end_time FROM incidents WHERE component_name=$1 AND start_time>=$2 AND end_time<=$3 ORDER BY id"

	rows, err := p.db.QueryContext(ctx, query, component, from, to)
	if err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, nil
}

func (p *PostgresRepo) CountIncidentsByService(ctx context.Context, service string, from, to time.Time) (int, error) {
	op := "CountIncidentsByService"
	query := "SELECT COUNT(*) FROM incidents WHERE service_name=$1 AND start_time>=$2 AND end_time<=$3"

	row := p.db.QueryRowContext(ctx, query, service, from, to)
	if err := row.Err(); err != nil {
		return 0, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, nil
}

func (p *PostgresRepo) CountIncidentsByComponent(ctx context.Context, component string, from, to time.Time) (int, error) {
	op := "CountIncidentsByComponent"
	query := "SELECT COUNT(*) FROM incidents WHERE component_name=$1 AND start_time>=$2 AND end_time<=$3"

	row := p.db.QueryRowContext(ctx, query, component, from, to)
	if err := row.Err(); err != nil {
		return 0, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, nil
}

// Report interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetReportByID(ctx context.Context, id int) (*models.Report, error) {
	op := "GetReportByID"
	query := "SELECT id, service_name, start_time, end_time FROM reports WHERE id=$1"

	row := p.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var r models.Report
	if err := row.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
		return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return &r, nil
}

func (p *PostgresRepo) ListReports(ctx context.Context, offset, limit int) ([]*models.Report, error) {
	op := "ListReports"
	query := "SELECT id, service_name, start_time, end_time FROM reports ORDER BY id OFFSET $1 LIMIT $2"

	rows, err := p.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	list := make([]*models.Report, 0)
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, nil
}

func (p *PostgresRepo) ListReportsByService(ctx context.Context, serviceName string, offset, limit int) ([]*models.Report, error) {
	op := "ListReportsByService"
	query := "SELECT id, service_name, start_time, end_time FROM %s WHERE service_name=$1 ORDER BY id OFFSET $2 LIMIT $3"

	rows, err := p.db.QueryContext(ctx, query, serviceName, offset, limit)
	if err != nil {
		return nil, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Report
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, nil
}

func (p *PostgresRepo) ListReportsByPeriod(ctx context.Context, from, to time.Time, offset, limit int) ([]*models.Report, error) {
	op := "ListReportsByPeriod"
	query := "SELECT id, service_name, start_time, end_time FROM reports WHERE start_time>=$1 AND end_time<=$2 ORDER BY id OFFSET $3 LIMIT $4"
	rows, err := p.db.QueryContext(ctx, query, from, to, offset, limit)
	if err != nil {
		return nil, &Error{Err: ErrNotFound, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Report
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, &Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, nil
}

func (p *PostgresRepo) CountReports(ctx context.Context) (int, error) {
	op := "CountReports"
	query := "SELECT COUNT(*) FROM reports"
	row := p.db.QueryRowContext(ctx, query)
	if err := row.Err(); err != nil {
		return 0, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, nil
}

func (p *PostgresRepo) CountReportsByService(ctx context.Context, serviceName string) (int, error) {
	op := "CountReportsByService"
	query := "SELECT COUNT(*) FROM reports WHERE service_name=$1"
	row := p.db.QueryRowContext(ctx, query, serviceName)
	if err := row.Err(); err != nil {
		return 0, &Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, &Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, nil
}
