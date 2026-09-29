package handlers

import (
	"github.com/finsight/backend/internal/middleware"
	"github.com/finsight/backend/internal/services"
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// InvoiceHandler manages AR invoices
type InvoiceHandler struct {
	invoiceRepo interface {
		ListInvoices(ctx interface{}, companyID interface{}, params interface{}) (interface{}, int64, error)
		GetInvoice(ctx interface{}, companyID, invoiceID interface{}) (interface{}, error)
		CreateInvoice(ctx interface{}, invoice interface{}) (interface{}, error)
		GetAgingBuckets(ctx interface{}, companyID, asOf interface{}) (interface{}, error)
		GetARStats(ctx interface{}, companyID, from, to interface{}) (interface{}, error)
	}
	fxSvc   *services.FXService
	auditRepo interface{}
	wcSvc   *services.WorkingCapitalService
	customerRepo interface{}
}

type CreateInvoiceRequest struct {
	CustomerID  string `json:"customer_id" binding:"required,uuid"`
	EntityID    string `json:"entity_id" binding:"required,uuid"`
	InvoiceDate string `json:"invoice_date" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
	Currency    string `json:"currency" binding:"required,len=3"`
	Reference   string `json:"reference"`
	Notes       string `json:"notes"`
	LineItems   []struct {
		Description string  `json:"description" binding:"required"`
		Quantity    float64 `json:"quantity" binding:"required,gt=0"`
		UnitPrice   string  `json:"unit_price" binding:"required"`
		TaxRate     string  `json:"tax_rate"`
	} `json:"line_items" binding:"required,min=1"`
}

type RecordPaymentRequest struct {
	PaymentDate   string `json:"payment_date" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
	PaymentMethod string `json:"payment_method"`
	Reference     string `json:"reference"`
	Notes         string `json:"notes"`
}

func NewInvoiceHandler(invoiceRepo, customerRepo, fxSvc, auditRepo, wcSvc interface{}) *InvoiceHandler {
	return &InvoiceHandler{}
}

// ListInvoices - GET /api/v1/invoices
func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	var params response.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "INVALID_PARAMS", err.Error())
		return
	}
	params.SetDefaults()

	companyID := middleware.GetCompanyID(c)
	_ = companyID

	// Placeholder response with correct structure
	response.OKWithMeta(c, []interface{}{}, &response.Meta{
		Page: params.Page, PageSize: params.PageSize, Total: 0,
	})
}

// GetInvoice - GET /api/v1/invoices/:id
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	_ = id
	response.NotFound(c, "INVOICE")
}

// CreateInvoice - POST /api/v1/invoices
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}
	response.NotFound(c, "INVOICE")
}

// UpdateInvoice - PUT /api/v1/invoices/:id
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	response.NotFound(c, "INVOICE")
}

// ApproveInvoice - POST /api/v1/invoices/:id/approve
func (h *InvoiceHandler) ApproveInvoice(c *gin.Context) {
	response.NotFound(c, "INVOICE")
}

// RecordPayment - POST /api/v1/invoices/:id/payment
func (h *InvoiceHandler) RecordPayment(c *gin.Context) {
	var req RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}
	response.NotFound(c, "INVOICE")
}

// DisputeInvoice - POST /api/v1/invoices/:id/dispute
func (h *InvoiceHandler) DisputeInvoice(c *gin.Context) {
	response.NotFound(c, "INVOICE")
}

// GetAging - GET /api/v1/invoices/aging
func (h *InvoiceHandler) GetAging(c *gin.Context) {
	response.OK(c, map[string]interface{}{
		"total_amount": 0,
		"buckets":      []interface{}{},
	})
}

// GetOverdue - GET /api/v1/invoices/overdue
func (h *InvoiceHandler) GetOverdue(c *gin.Context) {
	response.OKWithMeta(c, []interface{}{}, &response.Meta{})
}

// Export - GET /api/v1/invoices/export
func (h *InvoiceHandler) Export(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=invoices.csv")
	c.String(200, "invoice_number,customer,amount,status\n")
}

// GetARSummary - GET /api/v1/receivables/summary
func (h *InvoiceHandler) GetARSummary(c *gin.Context) {
	response.OK(c, map[string]interface{}{
		"total_outstanding": 0,
		"total_overdue":     0,
		"invoice_count":     0,
	})
}

// GetCollections - GET /api/v1/receivables/collections
func (h *InvoiceHandler) GetCollections(c *gin.Context) {
	response.OK(c, []interface{}{})
}
