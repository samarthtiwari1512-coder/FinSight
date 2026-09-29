package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	repo interface{}
}

func NewAuditHandler(repo interface{}) *AuditHandler {
	return &AuditHandler{repo}
}

func (h *AuditHandler) ListAuditLogs(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *AuditHandler) GetResourceHistory(c *gin.Context) { response.OK(c, []interface{}{}) }
