package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	repo interface{}
	fx interface{}
	audit interface{}
	cfg interface{}
}

func NewPaymentHandler(repo, fx, audit, cfg interface{}) *PaymentHandler {
	return &PaymentHandler{repo, fx, audit, cfg}
}

func (h *PaymentHandler) ListPayments(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *PaymentHandler) GetPayment(c *gin.Context) { response.NotFound(c, "PAYMENT") }
func (h *PaymentHandler) CreatePayment(c *gin.Context) { response.NotFound(c, "PAYMENT") }
func (h *PaymentHandler) ApprovePayment(c *gin.Context) { response.NotFound(c, "PAYMENT") }
func (h *PaymentHandler) RejectPayment(c *gin.Context) { response.NotFound(c, "PAYMENT") }
func (h *PaymentHandler) CancelPayment(c *gin.Context) { response.NotFound(c, "PAYMENT") }
