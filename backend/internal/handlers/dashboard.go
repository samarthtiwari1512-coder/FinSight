package handlers

import (
	"github.com/finsight/backend/internal/middleware"
	"github.com/finsight/backend/internal/services"
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// DashboardHandler serves the executive dashboard API
type DashboardHandler struct {
	dashSvc *services.DashboardService
}

func NewDashboardHandler(dashSvc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashSvc: dashSvc}
}

// GetDashboard returns full dashboard data
// GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	data, err := h.dashSvc.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, middleware.GetRequestID(c))
		return
	}

	response.OK(c, data)
}

// GetSummary returns the compact summary for header/widget use
// GET /api/v1/dashboard/summary
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	data, err := h.dashSvc.GetSummary(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, middleware.GetRequestID(c))
		return
	}

	response.OK(c, data)
}
