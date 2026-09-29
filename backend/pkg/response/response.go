package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Standard API response envelope
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Pagination params
type PaginationParams struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=200"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

func (p *PaginationParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 25
	}
	if p.SortDir == "" {
		p.SortDir = "desc"
	}
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// OK sends a 200 success response
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

// OKWithMeta sends a 200 success with pagination metadata
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data, Meta: meta})
}

// Created sends a 201 created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 error
func BadRequest(c *gin.Context, code, message string) {
	sendError(c, http.StatusBadRequest, code, message)
}

// Unauthorized sends a 401 error
func Unauthorized(c *gin.Context, code, message string) {
	sendError(c, http.StatusUnauthorized, code, message)
}

// Forbidden sends a 403 error
func Forbidden(c *gin.Context, code, message string) {
	sendError(c, http.StatusForbidden, code, message)
}

// NotFound sends a 404 error
func NotFound(c *gin.Context, resource string) {
	sendError(c, http.StatusNotFound, resource+"_NOT_FOUND",
		resource+" could not be found")
}

// Conflict sends a 409 error (e.g. duplicate)
func Conflict(c *gin.Context, code, message string) {
	sendError(c, http.StatusConflict, code, message)
}

// UnprocessableEntity sends a 422 error for validation failures
func UnprocessableEntity(c *gin.Context, code, message string) {
	sendError(c, http.StatusUnprocessableEntity, code, message)
}

// TooManyRequests sends a 429 rate limit error
func TooManyRequests(c *gin.Context) {
	sendError(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED",
		"Too many requests, please try again later")
}

// InternalError sends a 500 error — never exposing internal details
func InternalError(c *gin.Context, requestID string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &APIError{
			Code:      "INTERNAL_ERROR",
			Message:   "An unexpected error occurred. Please try again or contact support.",
			RequestID: requestID,
		},
	})
}

func sendError(c *gin.Context, status int, code, message string) {
	reqID := ""
	if v, exists := c.Get("request_id"); exists {
		reqID = v.(string)
	} else {
		reqID = uuid.New().String()
	}

	c.JSON(status, Response{
		Success: false,
		Error: &APIError{
			Code:      code,
			Message:   message,
			RequestID: reqID,
		},
	})
}

// BuildMeta creates pagination meta from params and total count
func BuildMeta(params PaginationParams, total int64) *Meta {
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}
	return &Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
