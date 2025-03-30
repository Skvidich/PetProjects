package models

import "time"

type Report struct {
	Name        string    `json:"name" db:"name"`
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`
}

type ComponentMetrics struct {
	Component string `json:"component" db:"component"`
	State     string `json:"state" db:"state"`
	Count     int    `json:"count" db:"count"`
}

type ServiceIncident struct {
	Name      string    `json:"name" db:"name"`
	Component string    `json:"component_name" db:"component_name"`
	State     string    `json:"status" db:"status"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}
