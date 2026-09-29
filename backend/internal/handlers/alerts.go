package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	svc interface{}
}

func NewAlertHandler(svc interface{}) *AlertHandler {
	return &AlertHandler{svc}
}

func (h *AlertHandler) ListAlerts(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *AlertHandler) GetAlert(c *gin.Context) { response.NotFound(c, "ALERT") }
func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) { response.NotFound(c, "ALERT") }
func (h *AlertHandler) DismissAlert(c *gin.Context) { response.NotFound(c, "ALERT") }
func (h *AlertHandler) GetAlertCounts(c *gin.Context) { response.OK(c, nil) }

func (h *AlertHandler) ListNotifications(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *AlertHandler) MarkRead(c *gin.Context) { response.OK(c, nil) }
func (h *AlertHandler) MarkAllRead(c *gin.Context) { response.OK(c, nil) }
func (h *AlertHandler) GetUnreadCount(c *gin.Context) { response.OK(c, map[string]int{"count": 0}) }
