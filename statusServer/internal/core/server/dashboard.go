package server

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *StatusServer) GetDashboardOverview(c *gin.Context) {
	op := "GetDashboardOverview"

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
	totalReports := len(reports)

	svcMap := make(map[string]struct{})
	for _, r := range reports {
		svcMap[r.ServiceName] = struct{}{}
	}
	totalIncidents := 0
	for svc := range svcMap {
		incs, err := s.repo.ListIncidentsByService(c.Request.Context(), svc, from, to)
		if err != nil {
			c.Error(NewError(err, op, "list incidents by service error", http.StatusInternalServerError))
			return
		}
		totalIncidents += len(incs)
	}

	c.JSON(http.StatusOK, gin.H{"overview": gin.H{
		"total_reports":   totalReports,
		"total_incidents": totalIncidents,
	}})
}

func (s *StatusServer) GetDashboardTopSlowest(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"top_slowest": []gin.H{}})
}

func (s *StatusServer) GetDashboardTopErrors(c *gin.Context) {
	op := "GetDashboardTopErrors"

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.Error(NewError(err, op, "invalid limit", http.StatusBadRequest))
		return
	}

	reports, err := s.repo.ListReportsByPeriod(c.Request.Context(), from, to, 0, 10000)
	if err != nil {
		c.Error(NewError(err, op, "list reports by period error", http.StatusInternalServerError))
		return
	}
	svcMap := make(map[string]struct{})
	for _, r := range reports {
		svcMap[r.ServiceName] = struct{}{}
	}

	type svcCount struct {
		Service string
		Count   int
	}
	var counts []svcCount
	for svc := range svcMap {
		incs, err := s.repo.ListIncidentsByService(c.Request.Context(), svc, from, to)
		if err != nil {
			c.Error(NewError(err, op, "list incidents by service error", http.StatusInternalServerError))
			return
		}
		counts = append(counts, svcCount{Service: svc, Count: len(incs)})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})
	if len(counts) > limit {
		counts = counts[:limit]
	}

	top := make([]gin.H, 0, len(counts))
	for _, sc := range counts {
		top = append(top, gin.H{"service": sc.Service, "incidents": sc.Count})
	}

	c.JSON(http.StatusOK, gin.H{"top_errors": top})
}

func (s *StatusServer) GetDashboardHeatmap(c *gin.Context) {
	op := "GetDashboardHeatmap"

	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}

	totalReps, err := s.repo.CountReports(c.Request.Context())
	if err != nil {
		c.Error(NewError(err, op, "count reports error", http.StatusInternalServerError))
		return
	}
	reps, err := s.repo.ListReports(c.Request.Context(), 0, totalReps)
	if err != nil {
		c.Error(NewError(err, op, "list reports error", http.StatusInternalServerError))
		return
	}
	svcMap := make(map[string]struct{})
	for _, r := range reps {
		svcMap[r.ServiceName] = struct{}{}
	}

	var heatmap []gin.H
	for svc := range svcMap {
		incs, err := s.repo.ListIncidentsByService(c.Request.Context(), svc, from, to)
		if err != nil {
			c.Error(NewError(err, op, "list incidents by service error", http.StatusInternalServerError))
			return
		}
		compCounts := make(map[string]int)
		for _, inc := range incs {
			compCounts[inc.ComponentName]++
		}
		for comp, cnt := range compCounts {
			heatmap = append(heatmap, gin.H{
				"service":   svc,
				"component": comp,
				"count":     cnt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"heatmap": heatmap})
}
