package storage

import (
	"dataCollector/pkg/types"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type SQLStorage struct {
	db             *sql.DB
	reportTable    string
	componentTable string
}

func NewSQLStorage(db *sql.DB, reportTable, componentTable string) (*SQLStorage, error) {
	return &SQLStorage{
		db:             db,
		reportTable:    reportTable,
		componentTable: componentTable,
	}, nil
}

func (s *SQLStorage) StoreRawReport(stat *types.ServiceStatus) error {

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Println("can't rollback transaction ", err.Error())
		}
	}()

	query := "INSERT INTO" + s.reportTable + " (service_name,start_time) VALUES ($1, $2) RETURNING id"
	result := tx.QueryRow(query, stat.Name, stat.Time.Format(time.DateTime))

	var reportId int64
	if err = result.Scan(&reportId); err != nil {
		return fmt.Errorf("can't read last inserted id %w", err)
	}

	query = "INSERT INTO" + s.componentTable + " (report_id,component_name,status) VALUES (" + strconv.FormatInt(reportId, 10) + ",$1,$2)"
	prep, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("can't prepare query %w", err)
	}
	defer func() {
		if err := prep.Close(); err != nil {
			log.Println("can't close statement ", err.Error())
		}
	}()

	for _, metric := range stat.Components {
		_, err = prep.Exec(metric.Name, metric.Status)
		if err != nil {
			return fmt.Errorf("can't execute query %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction %w", err)
	}

	return nil
}
