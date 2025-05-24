package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//
// Incidents processing
//

// ---------- Incidents Handlers ----------

func (s *StatusServer) GetListIncidents(c *gin.Context) {
	op := "GetListIncidents"
	svc := c.Query("service")
	comp := c.Query("component")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	incs, err := s.repo.ListIncidents(c.Request.Context(), svc, comp, from, to)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"incidents": incs})
}

func (s *StatusServer) GetIncidentByID(c *gin.Context) {
	op := "GetIncidentByID"
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(NewError(err, op, "invalid id", http.StatusBadRequest))
		return
	}
	inc, err := s.repo.GetIncidentByID(c.Request.Context(), id)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"incident": inc})
}

func (s *StatusServer) GetIncidentsByService(c *gin.Context) {
	op := "GetIncidentsByService"
	svc := c.Param("serviceName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	incs, err := s.repo.ListIncidentsByService(c.Request.Context(), svc, from, to)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"incidents": incs})
}

func (s *StatusServer) GetIncidentsByComponent(c *gin.Context) {
	op := "GetIncidentsByComponent"
	comp := c.Param("componentName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	incs, err := s.repo.ListIncidentsByComponent(c.Request.Context(), comp, from, to)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"incidents": incs})
}

func (s *StatusServer) GetIncidentsCountByService(c *gin.Context) {
	op := "GetIncidentsCountByService"
	svc := c.Param("serviceName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	cnt, err := s.repo.CountIncidentsByService(c.Request.Context(), svc, from, to)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": cnt})
}

func (s *StatusServer) GetIncidentsCountByComponent(c *gin.Context) {
	op := "GetIncidentsCountByComponent"
	comp := c.Param("componentName")
	from, ok := parseTimeParam(c, "from")
	if !ok {
		return
	}
	to, ok := parseTimeParam(c, "to")
	if !ok {
		return
	}
	cnt, err := s.repo.CountIncidentsByComponent(c.Request.Context(), comp, from, to)
	if err != nil {
		c.Error(NewError(err, op, "repo error", http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": cnt})
}
