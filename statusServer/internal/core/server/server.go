package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"statusServer/internal/core/midleware"
	"statusServer/internal/core/repository"
	"time"
)

type StatusServer struct {
	router *gin.Engine
	repo   repository.StatusRepository
}

func New(repo repository.StatusRepository) (*StatusServer, error) {

	var statServ StatusServer
	statServ.router = gin.New()
	statServ.router.Use(gin.Logger())
	statServ.router.Use(midleware.ErrHandler())
	statServ.repo = repo
	statServ.registerRoutes()
	return &statServ, nil
}

func (s *StatusServer) registerRoutes() {
	api := s.router.Group("/api")
	{

		reports := api.Group("/reports")
		{
			reports.GET("", s.GetListReports)
			reports.GET("/:id", s.GetReportByID)
			reports.GET("/service/:serviceName", s.GetReportsByService)
			reports.GET("/period", s.GetReportsByPeriod)
			reports.GET("/count", s.GetReportsCount)
			reports.GET("/service/:serviceName/count", s.GetReportsCountByService)
		}

		metrics := api.Group("/metrics")
		{
			metrics.GET("/:id", s.GetMetricByID)
			metrics.GET("/status/:status", s.GetMetricsByStatus)
		}

		api.GET("/reports/:reportId/metrics", s.GetMetricsByReport)
		api.GET("/reports/:reportId/metrics/status/:status", s.GetMetricsByReportAndStatus)
		api.GET("/reports/:reportId/metrics/count", s.GetMetricsCountByReport)
		// Components->Metrics
		api.GET("/components/:componentName/metrics", s.GetMetricsByComponent)

		incidents := api.Group("/incidents")
		{
			incidents.GET("", s.GetListIncidents)
			incidents.GET("/:id", s.GetIncidentByID)
			incidents.GET("/service/:serviceName", s.GetIncidentsByService)
			incidents.GET("/component/:componentName", s.GetIncidentsByComponent)
			incidents.GET("/service/:serviceName/count", s.GetIncidentsCountByService)
			incidents.GET("/component/:componentName/count", s.GetIncidentsCountByComponent)
		}

		services := api.Group("/services")
		{
			services.GET("", s.GetListServices)
			services.GET("/:serviceName", s.GetServiceSummary)
			services.GET("/:serviceName/uptime", s.GetServiceUptime)
			services.GET("/:serviceName/errors", s.GetServiceErrors)
		}

		components := api.Group("/components")
		{
			components.GET("", s.GetListComponents)
			components.GET("/service/:serviceName", s.GetComponentsByService)
			components.GET("/:componentName/summary", s.GetComponentSummary)
			components.GET("/:componentName/trend", s.GetComponentTrend)
		}

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/overview", s.GetDashboardOverview)
			dashboard.GET("/top-slowest", s.GetDashboardTopSlowest)
			dashboard.GET("/top-errors", s.GetDashboardTopErrors)
			dashboard.GET("/heatmap", s.GetDashboardHeatmap)
		}

		alerts := api.Group("/alerts")
		{
			alerts.GET("", s.GetListAlerts)

			alerts.GET("/history", s.GetAlertsHistory)
		}
	}
}

func (s *StatusServer) Start(addr string) error {
	return s.router.Run(addr)
}

func parseTimeParam(c *gin.Context, name string) (time.Time, bool) {
	str := c.Query(name)
	if str == "" {
		return time.Time{}, true
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		c.Error(NewError(err, "parseTimeParam", "invalid time format", http.StatusBadRequest))
		return time.Time{}, false
	}
	return t, true
}
