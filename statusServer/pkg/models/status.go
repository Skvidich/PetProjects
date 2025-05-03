package models

import "time"

type Report struct {
	ID          int       `json:"id,omitempty"`
	ServiceName string    `json:"service_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type ComponentMetric struct {
	ID            int    `json:"id,omitempty"`
	ReportID      int    `json:"report_id"`
	ComponentName string `json:"component_name"`
	Status        string `json:"status"`
	StatusCount   int16  `json:"status_count"`
}

type Incident struct {
	ID            int       `json:"id,omitempty"`
	ServiceName   string    `json:"service_name"`
	ComponentName string    `json:"component_name"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}
