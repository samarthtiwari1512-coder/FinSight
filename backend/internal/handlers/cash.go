package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type CashHandler struct {
	svc interface{}
	repo interface{}
	fx interface{}
}

func NewCashHandler(svc, repo, fx interface{}) *CashHandler {
	return &CashHandler{svc, repo, fx}
}

func (h *CashHandler) GetPosition(c *gin.Context) { response.OK(c, nil) }
func (h *CashHandler) GetPositionByBank(c *gin.Context) { response.OK(c, nil) }
func (h *CashHandler) GetPositionByCurrency(c *gin.Context) { response.OK(c, nil) }
func (h *CashHandler) GetPositionByEntity(c *gin.Context) { response.OK(c, nil) }
func (h *CashHandler) GetMovements(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *CashHandler) GetCashAccounts(c *gin.Context) { response.OK(c, []interface{}{}) }
