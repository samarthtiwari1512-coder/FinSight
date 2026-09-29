package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type ForecastHandler struct {
	svc interface{}
}

func NewForecastHandler(svc interface{}) *ForecastHandler {
	return &ForecastHandler{svc}
}

func (h *ForecastHandler) GetForecast(c *gin.Context) { response.OK(c, nil) }
func (h *ForecastHandler) GetForecast7d(c *gin.Context) { response.OK(c, nil) }
func (h *ForecastHandler) GetForecast30d(c *gin.Context) { response.OK(c, nil) }
func (h *ForecastHandler) GetForecast90d(c *gin.Context) { response.OK(c, nil) }
func (h *ForecastHandler) GetAccuracy(c *gin.Context) { response.OK(c, nil) }
func (h *ForecastHandler) Generate(c *gin.Context) { response.OK(c, nil) }
