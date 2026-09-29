package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type ScenarioHandler struct {
	svc interface{}
	repo interface{}
}

func NewScenarioHandler(svc, repo interface{}) *ScenarioHandler {
	return &ScenarioHandler{svc, repo}
}

func (h *ScenarioHandler) ListScenarios(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *ScenarioHandler) GetScenario(c *gin.Context) { response.NotFound(c, "SCENARIO") }
func (h *ScenarioHandler) CreateScenario(c *gin.Context) { response.NotFound(c, "SCENARIO") }
func (h *ScenarioHandler) RunScenario(c *gin.Context) { response.NotFound(c, "SCENARIO") }
func (h *ScenarioHandler) GetResults(c *gin.Context) { response.NotFound(c, "SCENARIO") }
func (h *ScenarioHandler) RunStressTest(c *gin.Context) { response.OK(c, nil) }
