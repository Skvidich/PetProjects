package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"statusServer/internal/core/repository"
	"strconv"
)

type StatusServer struct {
	router *gin.Engine
	repo   repository.StatusRepository
}

func New(repo repository.StatusRepository) (*StatusServer, error) {

	var statServ StatusServer
	statServ.router = gin.Default()
	statServ.repo = repo
	return &statServ, nil
}

func (s *StatusServer) Start(addr string) {
	go s.router.Run(addr)
}

//
// Reports processing
//

func (s *StatusServer) GetReportByID(c *gin.Context) {
	op := "GetReportByID"

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(NewError(err, op, "can't parse param", http.StatusBadRequest))
		return
	}

	report, err := s.repo.GetReportByID(c.Request.Context(), id)

	if err != nil {
		c.Error(NewError(err, op, "can't get data from repo", http.StatusInternalServerError))
		return
	}

	err = c.BindJSON(report)
	if err != nil {
		c.Error(NewError(err, op, "can't create JSON", http.StatusInternalServerError))
		return
	}

}

func (s *StatusServer) GetListReports(c *gin.Context) {
	op := "GetListReports"
	defOffset := "0"
	defLimit := "10"

	offset, err := strconv.Atoi(c.DefaultQuery("offset", defOffset))

	if err != nil {
		c.Error(NewError(err, op, "can't parse query", http.StatusBadRequest))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", defLimit))

	if err != nil {
		c.Error(NewError(err, op, "can't parse query", http.StatusBadRequest))
		return
	}

	reports, err := s.repo.ListReports(c.Request.Context(), offset, limit)

	if err != nil {
		c.Error(NewError(err, op, "can't get data from repo", http.StatusInternalServerError))
		return
	}

	err = c.BindJSON(&reports)
	if err != nil {
		c.Error(NewError(err, op, "can't create JSON", http.StatusInternalServerError))
		return
	}

}
