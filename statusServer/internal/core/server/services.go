package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *StatusServer) GetListServices(c *gin.Context) {
	op := "GetListServices"
	// get total reports to fetch all
	total, err := s.repo.CountReports(c.Request.Context())
	if err != nil {
		c.Error(NewError(err, op, "count reports error", http.StatusInternalServerError))
		return
	}
	// fetch all reports
	reps, err := s.repo.ListReports(c.Request.Context(), 0, total)
	if err != nil {
		c.Error(NewError(err, op, "list reports error", http.StatusInternalServerError))
		return
	}
	unique := make(map[string]struct{})
	for _, r := range reps {
		unique[r.ServiceName] = struct{}{}
	}
	services := make([]string, 0, len(unique))
	for svc := range unique {
		services = append(services, svc)
	}
	c.JSON(http.StatusOK, gin.H{"services": services})
}

func (s *StatusServer) GetServiceSummary(c *gin.Context) {
	op := "GetServiceSummary"
	service := c.Param("serviceName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	// count reports
	repCount, err := s.repo.CountReportsByService(c.Request.Context(), service)
	if err != nil {
		c.Error(NewError(err, op, "count reports error", http.StatusInternalServerError))
		return
	}
	// count incidents
	incCount, err := s.repo.CountIncidentsByService(c.Request.Context(), service, from, to)
	if err != nil {
		c.Error(NewError(err, op, "count incidents error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"report_count":   repCount,
			"incident_count": incCount,
		},
	})
}

func (s *StatusServer) GetServiceUptime(c *gin.Context) {
	op := "GetServiceUptime"
	service := c.Param("serviceName")

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}

	reports, err := s.repo.ListReportsByPeriod(c.Request.Context(), from, to, 0, 10000)
	if err != nil {
		c.Error(NewError(err, op, "list reports by period error", http.StatusInternalServerError))
		return
	}

	var totalCount, okCount int
	for _, r := range reports {
		if r.ServiceName != service {
			continue
		}

		metrics, err := s.repo.ListMetricsByReport(c.Request.Context(), r.ID)
		if err != nil {
			c.Error(NewError(err, op, "list metrics error", http.StatusInternalServerError))
			return
		}
		for _, m := range metrics {
			totalCount += int(m.StatusCount)
		}

		okMetrics, err := s.repo.ListMetricsByReportAndStatus(c.Request.Context(), r.ID, "OK")
		if err != nil {
			c.Error(NewError(err, op, "list OK metrics error", http.StatusInternalServerError))
			return
		}
		for _, m := range okMetrics {
			okCount += int(m.StatusCount)
		}
	}

	var uptimePct float64
	if totalCount > 0 {
		uptimePct = float64(okCount) / float64(totalCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{"uptime": uptimePct})
}

func (s *StatusServer) GetServiceErrors(c *gin.Context) {
	op := "GetServiceErrors"
	service := c.Param("serviceName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	incs, err := s.repo.ListIncidentsByService(c.Request.Context(), service, from, to)
	if err != nil {
		c.Error(NewError(err, op, "list incidents error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"errors": incs})
}
