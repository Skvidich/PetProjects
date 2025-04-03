package storage

import (
	"dataProcessor/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type SQLStorage struct {
	db             *sql.DB
	reportTable    string
	incidentTable  string
	componentTable string
}

func NewSQLStorage(db *sql.DB, reportTable, incidentTable, componentTable string) (*SQLStorage, error) {
	return &SQLStorage{
		db:             db,
		reportTable:    reportTable,
		incidentTable:  incidentTable,
		componentTable: componentTable,
	}, nil
}

func (s *SQLStorage) StoreReport(overview *models.Report, components []models.ComponentMetrics) error {

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {

		}
	}()

	query := "INSERT INTO " + s.reportTable + " (service_name,start_time,end_time) VALUES ($1, $2, $3) RETURNING id"
	result := tx.QueryRow(query, overview.Name, overview.PeriodStart.Format(time.DateTime), overview.PeriodEnd.Format(time.DateTime))

	var reportId int64
	if err = result.Scan(&reportId); err != nil {
		return fmt.Errorf("can't read last inserted id %w", err)
	}

	query = "INSERT INTO " + s.componentTable + " (report_id,component_name,status,status_count) VALUES (" + strconv.FormatInt(reportId, 10) + ",$1,$2,$3)"
	prep, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("can't prepare query %w", err)
	}
	defer func() {
		if err := prep.Close(); err != nil {

		}
	}()

	for _, metric := range components {
		_, err = prep.Exec(metric.Component, metric.State, metric.Count)
		if err != nil {
			return fmt.Errorf("can't execute query %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction %w", err)
	}

	return nil
}

func (s *SQLStorage) StoreIncident(incident *models.ServiceIncident) error {

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {

		}
	}()

	query := "INSERT INTO " + s.incidentTable + " (service_name,component_name,status,start_time,end_time) VALUES ($1, $2, $3,$4,$5)"
	_, err = tx.Exec(query, incident.Name, incident.Component, incident.State, incident.StartTime, incident.EndTime)
	if err != nil {
		return fmt.Errorf("can't execute query %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction %w", err)
	}

	return nil
}
