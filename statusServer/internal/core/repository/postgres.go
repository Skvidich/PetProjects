package repository

import (
	"context"
	"database/sql"
	"fmt"
	"statusServer/pkg/models"
	"time"
)

type PostgresRepo struct {
	// TODO add comprehensive logic
	db             *sql.DB
	ReportTable    string
	IncidentTable  string
	ComponentTable string
}

func NewPostgresRepo(db *sql.DB,
	ReportTable string,
	IncidentTable string,
	ComponentTable string) (StatusRepository, error) {
	return &PostgresRepo{
		db:             db,
		ReportTable:    ReportTable,
		IncidentTable:  IncidentTable,
		ComponentTable: ComponentTable,
	}, nil
}

func (p *PostgresRepo) Close() error {
	return nil
}

// Component interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetMetricByID(ctx context.Context, id int) (*models.ComponentMetric, Error) {

	op := "GetMetricsByID"
	query := fmt.Sprintf("SELECT id,report_id,component_name,status,status_count FROM %s WHERE id=$1", p.ComponentTable)

	res := p.db.QueryRowContext(ctx, query, id)
	if res.Err() != nil {
		return nil, Error{
			Err:     res.Err(),
			Problem: "can't execute query",
			Method:  op,
		}
	}

	var metric models.ComponentMetric

	err := res.Scan(&metric.ID, &metric.ReportID, &metric.ComponentName, &metric.Status, &metric.StatusCount)

	if err != nil {
		return nil, Error{
			Err:     err,
			Problem: "can't scan row",
			Method:  op,
		}
	}

	return &metric, Error{}

}

func (p *PostgresRepo) ListMetricsByReport(ctx context.Context, reportID int) ([]*models.ComponentMetric, Error) {

	op := "ListMetricsByReport"
	query := fmt.Sprintf("SELECT id,report_id,component_name,status,status_count FROM %s WHERE report_id=$1", p.ComponentTable)

	rows, err := p.db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, Error{
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
			return nil, Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, Error{}
}

func (p *PostgresRepo) ListMetricsByComponent(ctx context.Context, comp string) ([]*models.ComponentMetric, Error) {

	op := "ListMetricsByComponent"
	query := fmt.Sprintf("SELECT id,report_id,component_name,status,status_count FROM %s WHERE component_name=$1", p.ComponentTable)

	rows, err := p.db.QueryContext(ctx, query, comp)
	if err != nil {
		return nil, Error{
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
			return nil, Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, Error{}

}

func (p *PostgresRepo) ListMetricsByStatus(ctx context.Context, status string) ([]*models.ComponentMetric, Error) {
	op := "ListMetricsByStatus"
	query := fmt.Sprintf("SELECT id,report_id,component_name,status,status_count FROM %s WHERE status=$1", p.ComponentTable)

	rows, err := p.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, Error{
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
			return nil, Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, Error{}

}

func (p *PostgresRepo) ListMetricsByReportAndStatus(ctx context.Context, reportID int, status string) ([]*models.ComponentMetric, Error) {
	op := "ListMetricsByReportAndStatus"
	query := fmt.Sprintf("SELECT id,report_id,component_name,status,status_count FROM %s WHERE status=$1 AND report_id=$2", p.ComponentTable)

	rows, err := p.db.QueryContext(ctx, query, status, reportID)
	if err != nil {
		return nil, Error{
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
			return nil, Error{
				Err:     err,
				Problem: "can't scan row",
				Method:  op,
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, Error{}
}

func (p *PostgresRepo) CountMetricsByReport(ctx context.Context, reportID int) (int, Error) {
	op := "CountMetricsByReport"
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE report_id=$1", p.ComponentTable)

	res := p.db.QueryRowContext(ctx, query, reportID)
	if res.Err() != nil {
		return 0, Error{
			Err:     res.Err(),
			Problem: "can't execute query",
			Method:  op,
		}
	}

	var metricCount int

	err := res.Scan(&metricCount)

	if err != nil {
		return 0, Error{
			Err:     err,
			Problem: "can't scan row",
			Method:  op,
		}
	}

	return metricCount, Error{}
}

// Incident interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetIncidentByID(ctx context.Context, id int) (*models.Incident, Error) {
	op := "GetIncidentByID"
	query := fmt.Sprintf(
		"SELECT id, service_name, component_name, status, start_time, end_time FROM %s WHERE id=$1",
		p.IncidentTable,
	)
	row := p.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var inc models.Incident
	if err := row.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
		return nil, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return &inc, Error{}
}

func (p *PostgresRepo) ListIncidents(ctx context.Context, service, component string, from, to time.Time) ([]*models.Incident, Error) {
	op := "ListIncidents"
	query := fmt.Sprintf(
		"SELECT id, service_name, component_name, status, start_time, end_time FROM %s WHERE service_name=$1 AND component_name=$2 AND start_time>=$3 AND end_time<=$4 ORDER BY id",
		p.IncidentTable,
	)
	rows, err := p.db.QueryContext(ctx, query, service, component, from, to)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, Error{}
}

func (p *PostgresRepo) ListIncidentsByService(ctx context.Context, service string, from, to time.Time) ([]*models.Incident, Error) {
	op := "ListIncidentsByService"
	query := fmt.Sprintf(
		"SELECT id, service_name, component_name, status, start_time, end_time FROM %s WHERE service_name=$1 AND start_time>=$2 AND end_time<=$3 ORDER BY id",
		p.IncidentTable,
	)
	rows, err := p.db.QueryContext(ctx, query, service, from, to)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, Error{}
}

func (p *PostgresRepo) ListIncidentsByComponent(ctx context.Context, component string, from, to time.Time) ([]*models.Incident, Error) {
	op := "ListIncidentsByComponent"
	query := fmt.Sprintf(
		"SELECT id, service_name, component_name, status, start_time, end_time FROM %s WHERE component_name=$1 AND start_time>=$2 AND end_time<=$3 ORDER BY id",
		p.IncidentTable,
	)
	rows, err := p.db.QueryContext(ctx, query, component, from, to)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Incident
	for rows.Next() {
		var inc models.Incident
		if err = rows.Scan(&inc.ID, &inc.ServiceName, &inc.ComponentName, &inc.Status, &inc.StartTime, &inc.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &inc)
	}
	return list, Error{}
}

func (p *PostgresRepo) CountIncidentsByService(ctx context.Context, service string, from, to time.Time) (int, Error) {
	op := "CountIncidentsByService"
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE service_name=$1 AND start_time>=$2 AND end_time<=$3",
		p.IncidentTable,
	)
	row := p.db.QueryRowContext(ctx, query, service, from, to)
	if err := row.Err(); err != nil {
		return 0, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, Error{}
}

func (p *PostgresRepo) CountIncidentsByComponent(ctx context.Context, component string, from, to time.Time) (int, Error) {
	op := "CountIncidentsByComponent"
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE component_name=$1 AND start_time>=$2 AND end_time<=$3",
		p.IncidentTable,
	)
	row := p.db.QueryRowContext(ctx, query, component, from, to)
	if err := row.Err(); err != nil {
		return 0, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, Error{}
}

// Report interface
// -----------------------------------------------------------
//

func (p *PostgresRepo) GetReportByID(ctx context.Context, id int) (*models.Report, Error) {
	op := "GetReportByID"
	query := fmt.Sprintf(
		"SELECT id, service_name, start_time, end_time FROM %s WHERE id=$1",
		p.ReportTable,
	)
	row := p.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var r models.Report
	if err := row.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
		return nil, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return &r, Error{}
}

func (p *PostgresRepo) ListReports(ctx context.Context, offset, limit int) ([]*models.Report, Error) {
	op := "ListReports"
	query := fmt.Sprintf(
		"SELECT id, service_name, start_time, end_time FROM %s ORDER BY id OFFSET $1 LIMIT $2",
		p.ReportTable,
	)
	rows, err := p.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Report
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, Error{}
}

func (p *PostgresRepo) ListReportsByService(ctx context.Context, serviceName string, offset, limit int) ([]*models.Report, Error) {
	op := "ListReportsByService"
	query := fmt.Sprintf(
		"SELECT id, service_name, start_time, end_time FROM %s WHERE service_name=$1 ORDER BY id OFFSET $2 LIMIT $3",
		p.ReportTable,
	)
	rows, err := p.db.QueryContext(ctx, query, serviceName, offset, limit)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Report
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, Error{}
}

func (p *PostgresRepo) ListReportsByPeriod(ctx context.Context, from, to time.Time, offset, limit int) ([]*models.Report, Error) {
	op := "ListReportsByPeriod"
	query := fmt.Sprintf(
		"SELECT id, service_name, start_time, end_time FROM %s WHERE start_time>=$1 AND end_time<=$2 ORDER BY id OFFSET $3 LIMIT $4",
		p.ReportTable,
	)
	rows, err := p.db.QueryContext(ctx, query, from, to, offset, limit)
	if err != nil {
		return nil, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	defer rows.Close()
	var list []*models.Report
	for rows.Next() {
		var r models.Report
		if err = rows.Scan(&r.ID, &r.ServiceName, &r.StartTime, &r.EndTime); err != nil {
			return nil, Error{Err: err, Problem: "can't scan row", Method: op}
		}
		list = append(list, &r)
	}
	return list, Error{}
}

func (p *PostgresRepo) CountReports(ctx context.Context) (int, Error) {
	op := "CountReports"
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", p.ReportTable)
	row := p.db.QueryRowContext(ctx, query)
	if err := row.Err(); err != nil {
		return 0, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, Error{}
}

func (p *PostgresRepo) CountReportsByService(ctx context.Context, serviceName string) (int, Error) {
	op := "CountReportsByService"
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE service_name=$1", p.ReportTable)
	row := p.db.QueryRowContext(ctx, query, serviceName)
	if err := row.Err(); err != nil {
		return 0, Error{Err: err, Problem: "can't execute query", Method: op}
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, Error{Err: err, Problem: "can't scan row", Method: op}
	}
	return count, Error{}
}
