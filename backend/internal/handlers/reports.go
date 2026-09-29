package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	inv interface{}
	bill interface{}
	cash interface{}
	wc interface{}
	fx interface{}
}

func NewReportHandler(inv, bill, cash, wc, fx interface{}) *ReportHandler {
	return &ReportHandler{inv, bill, cash, wc, fx}
}

func (h *ReportHandler) ListReports(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *ReportHandler) GenerateCashPosition(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) GenerateWorkingCapital(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) GenerateARAgeing(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) GenerateAPAgeing(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) GenerateFXExposure(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) GenerateExecutiveSummary(c *gin.Context) { response.OK(c, nil) }
func (h *ReportHandler) DownloadReport(c *gin.Context) { response.NotFound(c, "REPORT") }
