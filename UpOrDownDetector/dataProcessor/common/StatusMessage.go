package common

import "time"

type StatusMessage struct {
	Name       string      `json:"name"`
	Time       time.Time   `json:"time"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func IsUpStatus(status string) bool {
	return false
}

func IsDownStatus(status string) bool {
	return false
}

type StatisticOverview struct {
	Name      string    `json:"name" db:"name"`             // Имя группы (аналог StatusMessage.Name)
	StartTime time.Time `json:"start_time" db:"start_time"` // Начало периода статистики
	EndTime   time.Time `json:"end_time" db:"end_time"`     // Конец периода статистики
}

// ComponentStatistic хранит статистику статусов для каждого компонента
type ComponentStatistic struct {
	ComponentName string `json:"component_name" db:"component_name"` // Имя компонента (Component.Name)
	Status        string `json:"status" db:"status"`                 // Статус компонента (Component.Status)
	Count         int    `json:"count" db:"count"`                   // Количество вхождений статуса
}

type Incident struct {
	Name          string    `json:"name" db:"name"`                     // Уникальный идентификатор
	ComponentName string    `json:"component_name" db:"component_name"` // Имя компонента (Component.Name)
	Status        string    `json:"status" db:"status"`                 // Статус (например, "down")
	StartTime     time.Time `json:"start_time" db:"start_time"`         // Время начала инцидента
	EndTime       time.Time `json:"end_time" db:"end_time"`             // Время окончания инцидента (если проблема решена) 	// Было ли отправлено оповещение
}
