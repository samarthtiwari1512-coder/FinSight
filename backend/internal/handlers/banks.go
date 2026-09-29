package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	repo interface{}
	audit interface{}
}

func NewBankHandler(repo, audit interface{}) *BankHandler {
	return &BankHandler{repo, audit}
}

func (h *BankHandler) ListBankAccounts(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *BankHandler) GetBankAccount(c *gin.Context) { response.NotFound(c, "BANK") }
func (h *BankHandler) CreateBankAccount(c *gin.Context) { response.NotFound(c, "BANK") }
func (h *BankHandler) UpdateBankAccount(c *gin.Context) { response.NotFound(c, "BANK") }
func (h *BankHandler) GetTransactions(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *BankHandler) ImportStatement(c *gin.Context) { response.NotFound(c, "BANK") }
