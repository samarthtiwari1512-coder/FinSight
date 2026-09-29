package handlers

import (
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo interface{}
	compRepo interface{}
	entRepo interface{}
	auditRepo interface{}
}

func NewAdminHandler(u, c, e, a interface{}) *AdminHandler {
	return &AdminHandler{u, c, e, a}
}

func (h *AdminHandler) ListUsers(c *gin.Context) { response.OKWithMeta(c, []interface{}{}, &response.Meta{}) }
func (h *AdminHandler) CreateUser(c *gin.Context) { response.NotFound(c, "USER") }
func (h *AdminHandler) UpdateUser(c *gin.Context) { response.NotFound(c, "USER") }
func (h *AdminHandler) AssignRole(c *gin.Context) { response.NotFound(c, "USER") }
func (h *AdminHandler) ListRoles(c *gin.Context) { response.OK(c, []interface{}{}) }
func (h *AdminHandler) GetSettings(c *gin.Context) { response.OK(c, nil) }
func (h *AdminHandler) UpdateSettings(c *gin.Context) { response.NotFound(c, "SETTINGS") }
func (h *AdminHandler) ListEntities(c *gin.Context) { response.OK(c, []interface{}{}) }
func (h *AdminHandler) CreateEntity(c *gin.Context) { response.NotFound(c, "ENTITY") }
