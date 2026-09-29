package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	repo interface{}
	invRepo interface{}
	audit interface{}
}

func NewCustomerHandler(repo, invRepo, audit interface{}) *CustomerHandler {
	return &CustomerHandler{repo, invRepo, audit}
}

func (h *CustomerHandler) ListCustomers(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *CustomerHandler) GetCustomer(c *gin.Context) { response.NotFound(c, "CUSTOMER") }
func (h *CustomerHandler) CreateCustomer(c *gin.Context) { response.NotFound(c, "CUSTOMER") }
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) { response.NotFound(c, "CUSTOMER") }
func (h *CustomerHandler) GetCustomerInvoices(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *CustomerHandler) GetCustomerRisk(c *gin.Context) { response.OK(c, nil) }
func (h *CustomerHandler) GetConcentration(c *gin.Context) { response.OK(c, []interface{}{}) }
