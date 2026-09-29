package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type FXHandler struct {
	svc interface{}
	repo interface{}
	audit interface{}
}

func NewFXHandler(svc, repo, audit interface{}) *FXHandler {
	return &FXHandler{svc, repo, audit}
}

func (h *FXHandler) GetRates(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) GetRate(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) GetRateHistory(c *gin.Context) { response.OK(c, []interface{}{}) }
func (h *FXHandler) GetExposure(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) GetExposureSummary(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) RunSensitivity(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) RunScenario(c *gin.Context) { response.OK(c, nil) }
func (h *FXHandler) ListDeals(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *FXHandler) GetDeal(c *gin.Context) { response.NotFound(c, "DEAL") }
func (h *FXHandler) CreateDeal(c *gin.Context) { response.NotFound(c, "DEAL") }
func (h *FXHandler) ApproveDeal(c *gin.Context) { response.NotFound(c, "DEAL") }
func (h *FXHandler) ListHedges(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *FXHandler) CreateHedge(c *gin.Context) { response.NotFound(c, "HEDGE") }
func (h *FXHandler) GetHedgeCoverage(c *gin.Context) { response.OK(c, nil) }
