package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type RiskHandler struct {
	svc interface{}
	repo interface{}
}

func NewRiskHandler(svc, repo interface{}) *RiskHandler {
	return &RiskHandler{svc, repo}
}

func (h *RiskHandler) ListEvents(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *RiskHandler) GetEvent(c *gin.Context) { response.NotFound(c, "EVENT") }
func (h *RiskHandler) AcknowledgeEvent(c *gin.Context) { response.NotFound(c, "EVENT") }
func (h *RiskHandler) ResolveEvent(c *gin.Context) { response.NotFound(c, "EVENT") }
func (h *RiskHandler) ListLimits(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *RiskHandler) CreateLimit(c *gin.Context) { response.NotFound(c, "LIMIT") }
func (h *RiskHandler) UpdateLimit(c *gin.Context) { response.NotFound(c, "LIMIT") }
func (h *RiskHandler) GetRiskSummary(c *gin.Context) { response.OK(c, nil) }
