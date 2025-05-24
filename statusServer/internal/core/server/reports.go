package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//
// Reports processing
//

// ---------- Reports Handlers ----------

func (s *StatusServer) GetListReports(c *gin.Context) {
	op := "GetListReports"
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(NewError(err, op, "invalid offset", http.StatusBadRequest))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.Error(NewError(err, op, "invalid limit", http.StatusBadRequest))
		return
	}
	reps, err := s.repo.ListReports(c.Request.Context(), offset, limit)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reps})
}

func (s *StatusServer) GetReportByID(c *gin.Context) {
	op := "GetReportByID"
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(NewError(err, op, "invalid id", http.StatusBadRequest))
		return
	}
	rep, err := s.repo.GetReportByID(c.Request.Context(), id)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": rep})
}

func (s *StatusServer) GetReportsByService(c *gin.Context) {
	op := "GetReportsByService"
	service := c.Param("serviceName")
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(NewError(err, op, "invalid offset", http.StatusBadRequest))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.Error(NewError(err, op, "invalid limit", http.StatusBadRequest))
		return
	}
	reps, err := s.repo.ListReportsByService(c.Request.Context(), service, offset, limit)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reps})
}

func (s *StatusServer) GetReportsByPeriod(c *gin.Context) {
	op := "GetReportsByPeriod"
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(NewError(err, op, "invalid offset", http.StatusBadRequest))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.Error(NewError(err, op, "invalid limit", http.StatusBadRequest))
		return
	}
	reps, err := s.repo.ListReportsByPeriod(c.Request.Context(), from, to, offset, limit)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reps})
}

func (s *StatusServer) GetReportsCount(c *gin.Context) {
	op := "GetReportsCount"
	count, err := s.repo.CountReports(c.Request.Context())
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (s *StatusServer) GetReportsCountByService(c *gin.Context) {
	op := "GetReportsCountByService"
	service := c.Param("serviceName")
	count, err := s.repo.CountReportsByService(c.Request.Context(), service)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
