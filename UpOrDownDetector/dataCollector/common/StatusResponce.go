package common

import "time"

type StatusResponse struct {
	Name       string      `json:"name"`
	Time       time.Time   `json:"time"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
