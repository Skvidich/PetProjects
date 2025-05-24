package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//
// Metrics processing
//

// ---------- Metrics Handlers ----------

func (s *StatusServer) GetMetricByID(c *gin.Context) {
	op := "GetMetricByID"
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(NewError(err, op, "invalid id", http.StatusBadRequest))
		return
	}
	m, err := s.repo.GetMetricByID(c.Request.Context(), id)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"metric": m})
}

func (s *StatusServer) GetMetricsByReport(c *gin.Context) {
	op := "GetMetricsByReport"
	rid, err := strconv.Atoi(c.Param("reportId"))
	if err != nil {
		c.Error(NewError(err, op, "invalid reportId", http.StatusBadRequest))
		return
	}
	ms, err := s.repo.ListMetricsByReport(c.Request.Context(), rid)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": ms})
}

func (s *StatusServer) GetMetricsByComponent(c *gin.Context) {
	op := "GetMetricsByComponent"
	comp := c.Param("componentName")
	ms, err := s.repo.ListMetricsByComponent(c.Request.Context(), comp)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": ms})
}

func (s *StatusServer) GetMetricsByStatus(c *gin.Context) {
	op := "GetMetricsByStatus"
	status := c.Param("status")
	ms, err := s.repo.ListMetricsByStatus(c.Request.Context(), status)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": ms})
}

func (s *StatusServer) GetMetricsByReportAndStatus(c *gin.Context) {
	op := "GetMetricsByReportAndStatus"
	rid, err := strconv.Atoi(c.Param("reportId"))
	if err != nil {
		c.Error(NewError(err, op, "invalid reportId", http.StatusBadRequest))
		return
	}
	status := c.Param("status")
	ms, err := s.repo.ListMetricsByReportAndStatus(c.Request.Context(), rid, status)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": ms})
}

func (s *StatusServer) GetMetricsCountByReport(c *gin.Context) {
	op := "GetMetricsCountByReport"
	rid, err := strconv.Atoi(c.Param("reportId"))
	if err != nil {
		c.Error(NewError(err, op, "invalid reportId", http.StatusBadRequest))
		return
	}
	cnt, err := s.repo.CountMetricsByReport(c.Request.Context(), rid)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": cnt})
}
