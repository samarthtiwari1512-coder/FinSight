package handlers

import (
	"net/http"

	"github.com/finsight/backend/internal/auth"
	"github.com/finsight/backend/internal/middleware"
	"github.com/finsight/backend/internal/repositories"
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc   *auth.Service
	auditRepo repositories.AuditRepository
}

func NewAuthHandler(authSvc *auth.Service, auditRepo repositories.AuditRepository) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, auditRepo: auditRepo}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login - POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Email and password are required")
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	result, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password, ip, userAgent)
	if err != nil {
		switch err.Error() {
		case auth.ErrInvalidCredentials:
			response.Unauthorized(c, "INVALID_CREDENTIALS", "Invalid email or password")
		case auth.ErrAccountLocked:
			response.Unauthorized(c, "ACCOUNT_LOCKED", "Account is locked. Please contact your administrator.")
		case auth.ErrAccountInactive:
			response.Unauthorized(c, "ACCOUNT_INACTIVE", "Account is inactive. Please contact your administrator.")
		default:
			response.InternalError(c, middleware.GetRequestID(c))
		}
		return
	}

	response.OK(c, result)
}

// Refresh - POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Refresh token is required")
		return
	}

	ip := c.ClientIP()
	result, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken, ip)
	if err != nil {
		response.Unauthorized(c, "INVALID_REFRESH_TOKEN", "Refresh token is invalid or expired")
		return
	}

	response.OK(c, result)
}

// Logout - POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	c.ShouldBindJSON(&req) // Optional body

	userID := middleware.GetUserID(c)
	ip := c.ClientIP()

	if req.RefreshToken != "" {
		h.authSvc.Logout(c.Request.Context(), req.RefreshToken, userID, ip)
	}

	response.OK(c, gin.H{"message": "Logged out successfully"})
}

// Me - GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "USER")
		return
	}

	response.OK(c, user)
}
