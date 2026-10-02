package handlers

import (
	"context"
	"time"

	"github.com/finsight/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// HealthHandler provides system health endpoints
type HealthHandler struct {
	db    interface{ Ping(context.Context) error }
	redis interface{ Ping(ctx context.Context) *redis.StatusCmd }
	fxSvc interface{}
}

func NewHealthHandler(db interface{ Ping(context.Context) error }, rdb *redis.Client, fxSvc interface{}) *HealthHandler {
	return &HealthHandler{db: db, redis: rdb, fxSvc: fxSvc}
}

type HealthCheck struct {
	Status     string                 `json:"status"`
	Version    string                 `json:"version"`
	Timestamp  time.Time              `json:"timestamp"`
	Components map[string]interface{} `json:"components,omitempty"`
}

// Health - GET /health (public, lightweight)
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, HealthCheck{
		Status:    "ok",
		Version:   "1.0.0",
		Timestamp: time.Now(),
	})
}

// Ready - GET /ready (readiness probe, checks dependencies)
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx := c.Request.Context()

	// Check DB
	dbOK := true
	if err := h.db.Ping(ctx); err != nil {
		dbOK = false
	}

	// Check Redis
	redisOK := true
	if rc, ok := h.redis.(*redis.Client); ok {
		if err := rc.Ping(ctx).Err(); err != nil {
			redisOK = false
		}
	}

	if !dbOK {
		c.JSON(503, HealthCheck{
			Status:    "not_ready",
			Timestamp: time.Now(),
			Components: map[string]interface{}{
				"database": map[string]bool{"ok": dbOK},
				"redis":    map[string]bool{"ok": redisOK},
			},
		})
		return
	}

	c.JSON(200, HealthCheck{
		Status:    "ready",
		Timestamp: time.Now(),
		Components: map[string]interface{}{
			"database": map[string]bool{"ok": dbOK},
			"redis":    map[string]bool{"ok": redisOK},
		},
	})
}

// DetailedHealth - GET /api/v1/system/health (admin only)
func (h *HealthHandler) DetailedHealth(c *gin.Context) {
	ctx := c.Request.Context()

	dbOK := h.db.Ping(ctx) == nil
	var redisOK bool
	if rc, ok := h.redis.(*redis.Client); ok {
		redisOK = rc.Ping(ctx).Err() == nil
	}

	status := "healthy"
	if !dbOK {
		status = "degraded"
	}

	response.OK(c, map[string]interface{}{
		"status":    status,
		"timestamp": time.Now(),
		"database":  map[string]bool{"connected": dbOK},
		"redis":     map[string]bool{"connected": redisOK},
		"version":   "1.0.0",
		"env":       "demo",
	})
}
