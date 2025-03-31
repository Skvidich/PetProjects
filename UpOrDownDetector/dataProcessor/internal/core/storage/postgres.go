package storage

import (
	"dataProcessor/pkg/models"
	"fmt"
	"time"
)

type SQLStore struct {
}

func NewSQLStore() *SQLStore {
	return &SQLStore{}
}

func (s *SQLStore) StoreReport(report *models.Report, metrics []models.ComponentMetrics) error {
	fmt.Println("statistic saved")
	fmt.Println(report.Name, " ", report.PeriodStart.String(), " ", report.PeriodEnd.String())
	for _, m := range metrics {
		fmt.Printf("Metric: %s - %s (%d)\n", m.Component, m.State, m.Count)
	}
	return nil
}

func (s *SQLStore) StoreIncident(incident *models.ServiceIncident) error {
	fmt.Printf("Recording incident: %s [%s] %s-%s\n",
		incident.Name,
		incident.Component,
		incident.StartTime.Format(time.RFC1123),
		incident.EndTime.Format(time.RFC1123))
	return nil
}
