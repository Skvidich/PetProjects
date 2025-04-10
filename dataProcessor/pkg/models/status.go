package models

import "time"

type ServiceStatus struct {
	Name       string      `json:"name"`
	Time       time.Time   `json:"time"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func IsUp(status string) bool {
	return status == "operational"
}

func IsDown(status string) bool {
	return status != "operational"
}
