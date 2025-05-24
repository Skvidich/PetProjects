package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *StatusServer) GetListAlerts(c *gin.Context) {
	op := "GetListAlerts"

	from := time.Time{}
	to := time.Now().UTC()
	incs, err := s.repo.ListIncidents(c.Request.Context(), "", "", from, to)
	if err != nil {
		c.Error(NewError(err, op, "list incidents error", http.StatusInternalServerError))
		return
	}

	var alerts []gin.H
	for _, inc := range incs {
		if inc.EndTime.IsZero() {
			alerts = append(alerts, gin.H{
				"id":         inc.ID,
				"service":    inc.ServiceName,
				"component":  inc.ComponentName,
				"status":     inc.Status,
				"start_time": inc.StartTime,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (s *StatusServer) GetAlertsHistory(c *gin.Context) {
	op := "GetAlertsHistory"

	service := c.Query("service")
	component := c.Query("component")

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}

	incs, err := s.repo.ListIncidents(c.Request.Context(), service, component, from, to)
	if err != nil {
		c.Error(NewError(err, op, "list incidents error", http.StatusInternalServerError))
		return
	}

	history := make([]gin.H, 0, len(incs))
	for _, inc := range incs {
		history = append(history, gin.H{
			"id":         inc.ID,
			"service":    inc.ServiceName,
			"component":  inc.ComponentName,
			"status":     inc.Status,
			"start_time": inc.StartTime,
			"end_time":   inc.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
