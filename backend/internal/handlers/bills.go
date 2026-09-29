package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	repo interface{}
	suppRepo interface{}
	fx interface{}
	audit interface{}
}

func NewBillHandler(repo, suppRepo, fx, audit interface{}) *BillHandler {
	return &BillHandler{repo, suppRepo, fx, audit}
}

func (h *BillHandler) ListBills(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *BillHandler) GetBill(c *gin.Context) { response.NotFound(c, "BILL") }
func (h *BillHandler) CreateBill(c *gin.Context) { response.NotFound(c, "BILL") }
func (h *BillHandler) UpdateBill(c *gin.Context) { response.NotFound(c, "BILL") }
func (h *BillHandler) ApproveBill(c *gin.Context) { response.NotFound(c, "BILL") }
func (h *BillHandler) RecordPayment(c *gin.Context) { response.NotFound(c, "BILL") }
func (h *BillHandler) GetAging(c *gin.Context) { response.OK(c, map[string]interface{}{"buckets": []interface{}{}}) }
func (h *BillHandler) GetOverdue(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *BillHandler) GetAPSummary(c *gin.Context) { response.OK(c, map[string]interface{}{}) }
func (h *BillHandler) GetDueSoon(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
