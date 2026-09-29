package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type WorkingCapitalHandler struct {
	svc interface{}
}

func NewWorkingCapitalHandler(svc interface{}) *WorkingCapitalHandler {
	return &WorkingCapitalHandler{svc}
}

func (h *WorkingCapitalHandler) GetKPIs(c *gin.Context) { response.OK(c, nil) }
func (h *WorkingCapitalHandler) GetTrend(c *gin.Context) { response.OK(c, []interface{}{}) }
func (h *WorkingCapitalHandler) GetOpportunities(c *gin.Context) { response.OK(c, []interface{}{}) }
func (h *WorkingCapitalHandler) GetCCC(c *gin.Context) { response.OK(c, nil) }
