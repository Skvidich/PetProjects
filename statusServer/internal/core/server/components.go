package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *StatusServer) GetListComponents(c *gin.Context) {
	op := "GetListComponents"
	totalReports, err := s.repo.CountReports(c.Request.Context())
	if err != nil {
		c.Error(NewError(err, op, "count reports error", http.StatusInternalServerError))
		return
	}
	reps, err := s.repo.ListReports(c.Request.Context(), 0, totalReports)
	if err != nil {
		c.Error(NewError(err, op, "list reports error", http.StatusInternalServerError))
		return
	}
	unique := make(map[string]struct{})
	for _, r := range reps {
		mets, err := s.repo.ListMetricsByReport(c.Request.Context(), r.ID)
		if err != nil {
			c.Error(NewError(err, op, "list metrics error", http.StatusInternalServerError))
			return
		}
		for _, m := range mets {
			unique[m.ComponentName] = struct{}{}
		}
	}
	components := make([]string, 0, len(unique))
	for comp := range unique {
		components = append(components, comp)
	}
	c.JSON(http.StatusOK, gin.H{"components": components})
}

func (s *StatusServer) GetComponentsByService(c *gin.Context) {
	op := "GetComponentsByService"
	service := c.Param("serviceName")

	totalReports, err := s.repo.CountReportsByService(c.Request.Context(), service)
	if err != nil {
		c.Error(NewError(err, op, "count reports by service error", http.StatusInternalServerError))
		return
	}
	reps, err := s.repo.ListReportsByService(c.Request.Context(), service, 0, totalReports)
	if err != nil {
		c.Error(NewError(err, op, "list reports by service error", http.StatusInternalServerError))
		return
	}

	unique := make(map[string]struct{})
	for _, r := range reps {
		mets, err := s.repo.ListMetricsByReport(c.Request.Context(), r.ID)
		if err != nil {
			c.Error(NewError(err, op, "list metrics error", http.StatusInternalServerError))
			return
		}
		for _, m := range mets {
			unique[m.ComponentName] = struct{}{}
		}
	}
	components := make([]string, 0, len(unique))
	for comp := range unique {
		components = append(components, comp)
	}
	c.JSON(http.StatusOK, gin.H{"components": components})
}

func (s *StatusServer) GetComponentSummary(c *gin.Context) {
	op := "GetComponentSummary"
	component := c.Param("componentName")

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}

	mets, err := s.repo.ListMetricsByComponent(c.Request.Context(), component)
	if err != nil {
		c.Error(NewError(err, op, "list metrics by component error", http.StatusInternalServerError))
		return
	}
	var totalMetrics int
	for _, m := range mets {
		totalMetrics += int(m.StatusCount)
	}

	incCount, err := s.repo.CountIncidentsByComponent(c.Request.Context(), component, from, to)
	if err != nil {
		c.Error(NewError(err, op, "count incidents by component error", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"metric_count":   totalMetrics,
			"incident_count": incCount,
		},
	})
}

func (s *StatusServer) GetComponentTrend(c *gin.Context) {
	op := "GetComponentTrend"
	component := c.Param("componentName")

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}

	reps, err := s.repo.ListReportsByPeriod(c.Request.Context(), from, to, 0, 10000)
	if err != nil {
		c.Error(NewError(err, op, "list reports by period error", http.StatusInternalServerError))
		return
	}

	trend := make([]gin.H, 0, len(reps))
	for _, r := range reps {
		mets, err := s.repo.ListMetricsByReport(c.Request.Context(), r.ID)
		if err != nil {
			c.Error(NewError(err, op, "list metrics error", http.StatusInternalServerError))
			return
		}
		sum := 0
		for _, m := range mets {
			if m.ComponentName == component {
				sum += int(m.StatusCount)
			}
		}
		trend = append(trend, gin.H{
			"report_id":    r.ID,
			"time":         r.StartTime,
			"metric_count": sum,
		})
	}
	c.JSON(http.StatusOK, gin.H{"trend": trend})
}
