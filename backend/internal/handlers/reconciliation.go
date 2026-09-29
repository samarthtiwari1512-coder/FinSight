package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReconciliationHandler struct {
	svc interface{}
	repo interface{}
}

func NewReconciliationHandler(svc, repo interface{}) *ReconciliationHandler {
	return &ReconciliationHandler{svc, repo}
}

func (h *ReconciliationHandler) ListReconciliations(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *ReconciliationHandler) GetDashboard(c *gin.Context) { response.OK(c, nil) }
func (h *ReconciliationHandler) ManualMatch(c *gin.Context) { response.OK(c, nil) }
func (h *ReconciliationHandler) GetUnmatched(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *ReconciliationHandler) GetExceptions(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
