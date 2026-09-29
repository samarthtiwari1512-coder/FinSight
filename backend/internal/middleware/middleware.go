package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/finsight/backend/internal/auth"
	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	UserIDKey    = "user_id"
	CompanyIDKey = "company_id"
	UserEmailKey = "user_email"
	UserRolesKey = "user_roles"
	ClaimsKey    = "claims"
	RequestIDKey = "request_id"
)

// RequestID injects a unique request ID into every request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger logs structured request/response info
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		reqID, _ := c.Get(RequestIDKey)
		userID, _ := c.Get(UserIDKey)

		level := zap.InfoLevel
		if status >= 500 {
			level = zap.ErrorLevel
		} else if status >= 400 {
			level = zap.WarnLevel
		}

		fields := []zap.Field{
			zap.String("request_id", fmt.Sprintf("%v", reqID)),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if userID != nil {
			fields = append(fields, zap.Any("user_id", userID))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		logger.Log(level, "HTTP Request", fields...)
	}
}

// Auth validates JWT and sets user context
func Auth(authSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c, "MISSING_TOKEN", "Authorization token is required")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := authSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "INVALID_TOKEN", "Token is invalid or expired")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(CompanyIDKey, claims.CompanyID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRolesKey, claims.Roles)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// RequirePermission checks if the authenticated user has a specific permission
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get(UserRolesKey)
		if !exists {
			response.Forbidden(c, "PERMISSION_DENIED", "Insufficient permissions")
			c.Abort()
			return
		}

		userRoles := roles.([]string)

		// Admin has all permissions
		for _, role := range userRoles {
			if role == "ADMIN" || role == "CFO" {
				c.Next()
				return
			}
		}

		// Check specific permissions via role-permission mapping
		if !hasPermission(userRoles, resource, action) {
			response.Forbidden(c, "PERMISSION_DENIED",
				"You do not have permission to "+action+" "+resource)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole ensures the user has at least one of the given roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get(UserRolesKey)
		if !exists {
			response.Forbidden(c, "ROLE_REQUIRED", "Role check failed")
			c.Abort()
			return
		}

		actualRoles := userRoles.([]string)
		for _, required := range roles {
			for _, actual := range actualRoles {
				if actual == required || actual == "ADMIN" {
					c.Next()
					return
				}
			}
		}

		response.Forbidden(c, "INSUFFICIENT_ROLE", "This action requires elevated privileges")
		c.Abort()
	}
}

// SecureHeaders adds security headers to all responses
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

// GetUserID extracts the user ID from context
func GetUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(UserIDKey)
	if id, ok := v.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// GetCompanyID extracts the company ID from context
func GetCompanyID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(CompanyIDKey)
	if id, ok := v.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// GetRequestID extracts the request ID from context
func GetRequestID(c *gin.Context) string {
	v, _ := c.Get(RequestIDKey)
	if id, ok := v.(string); ok {
		return id
	}
	return ""
}

// Permission mapping — in production, this would be loaded from DB
// This is a simplified static mapping for demonstration
var rolePermissions = map[string][]string{
	"CFO": {
		"dashboard:read", "cash:read", "cash:export",
		"bank_accounts:read", "customers:read", "suppliers:read",
		"invoices:read", "invoices:export", "invoices:approve",
		"bills:read", "bills:approve",
		"payments:read", "payments:approve",
		"fx:read", "fx:approve", "forecast:read",
		"risk:read", "risk:manage", "scenarios:read", "scenarios:create",
		"reconciliation:read", "alerts:read", "alerts:manage",
		"reports:read", "reports:create", "audit_logs:read",
		"working_capital:read",
	},
	"TREASURY_MANAGER": {
		"dashboard:read", "cash:read", "cash:export",
		"bank_accounts:read", "bank_accounts:create", "bank_accounts:update",
		"fx:read", "fx:create", "fx:approve",
		"forecast:read", "forecast:create",
		"scenarios:read", "scenarios:create",
		"alerts:read", "reports:read", "reports:create",
		"reconciliation:read", "reconciliation:create",
		"working_capital:read", "risk:read",
	},
	"FINANCE_MANAGER": {
		"dashboard:read", "customers:read", "customers:create", "customers:update",
		"suppliers:read", "suppliers:create", "suppliers:update",
		"invoices:read", "invoices:create", "invoices:update", "invoices:approve", "invoices:export",
		"bills:read", "bills:create", "bills:update", "bills:approve",
		"payments:read", "payments:create", "payments:approve",
		"reconciliation:read", "reconciliation:create",
		"working_capital:read", "reports:read", "reports:create",
		"alerts:read", "cash:read",
	},
	"AR_USER": {
		"customers:read", "customers:create", "customers:update",
		"invoices:read", "invoices:create", "invoices:update", "invoices:export",
		"alerts:read", "reports:read", "working_capital:read",
	},
	"AP_USER": {
		"suppliers:read", "suppliers:create", "suppliers:update",
		"bills:read", "bills:create", "bills:update",
		"payments:read", "payments:create",
		"alerts:read", "reports:read",
	},
	"RISK_ANALYST": {
		"dashboard:read", "fx:read", "risk:read", "risk:manage",
		"scenarios:read", "scenarios:create",
		"alerts:read", "reports:read", "working_capital:read",
	},
	"AUDITOR": {
		"dashboard:read", "cash:read", "bank_accounts:read",
		"customers:read", "suppliers:read", "invoices:read",
		"bills:read", "payments:read", "fx:read",
		"reconciliation:read", "audit_logs:read", "alerts:read",
		"reports:read", "working_capital:read",
	},
	"ADMIN": {"*:*"}, // Admin has everything
}

func hasPermission(userRoles []string, resource, action string) bool {
	key := resource + ":" + action
	for _, role := range userRoles {
		perms, ok := rolePermissions[role]
		if !ok {
			continue
		}
		for _, p := range perms {
			if p == "*:*" || p == key || p == resource+":*" {
				return true
			}
		}
	}
	return false
}


