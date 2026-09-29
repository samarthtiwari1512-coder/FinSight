package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	repo interface{}
	billRepo interface{}
	audit interface{}
}

func NewSupplierHandler(repo, billRepo, audit interface{}) *SupplierHandler {
	return &SupplierHandler{repo, billRepo, audit}
}

func (h *SupplierHandler) ListSuppliers(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *SupplierHandler) GetSupplier(c *gin.Context) { response.NotFound(c, "SUPPLIER") }
func (h *SupplierHandler) CreateSupplier(c *gin.Context) { response.NotFound(c, "SUPPLIER") }
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) { response.NotFound(c, "SUPPLIER") }
func (h *SupplierHandler) GetSupplierBills(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *SupplierHandler) GetConcentration(c *gin.Context) { response.OK(c, []interface{}{}) }
