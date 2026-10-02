package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/finsight/backend/internal/config"
	"github.com/finsight/backend/internal/models"
	"github.com/finsight/backend/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountLocked      = errors.New("account temporarily locked due to too many failed attempts")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrTokenInvalid       = errors.New("token is invalid or expired")
	ErrTokenRevoked       = errors.New("token has been revoked")
)

const (
	maxFailedAttempts = 5
	lockDuration      = 15 * time.Minute
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	TokenType string    `json:"token_type"` // access, refresh
	jwt.RegisteredClaims
}

type Service struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
	auditRepo repositories.AuditRepository
	cfg       *config.JWTConfig
}

func NewService(
	userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	auditRepo repositories.AuditRepository,
	cfg *config.JWTConfig,
) *Service {
	return &Service{userRepo: userRepo, tokenRepo: tokenRepo, auditRepo: auditRepo, cfg: cfg}
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"` // seconds
	User         *models.User `json:"user"`
}

func (s *Service) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Log failed attempt (don't reveal whether email exists)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	if user.IsLocked() {
		return nil, ErrAccountLocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Increment failed attempts
		s.userRepo.IncrementFailedAttempts(ctx, user.ID, maxFailedAttempts, lockDuration)
		return nil, ErrInvalidCredentials
	}

	// Reset failed attempts on success
	s.userRepo.ResetFailedAttempts(ctx, user.ID)
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Load roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("loading roles: %w", err)
	}
	user.Roles = roles

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	rawRefresh, tokenHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Store refresh token
	expiresAt := time.Now().Add(time.Duration(s.cfg.RefreshExpiryDays) * 24 * time.Hour)
	err = s.tokenRepo.StoreRefreshToken(ctx, user.ID, tokenHash, expiresAt, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("storing refresh token: %w", err)
	}

	// Audit log
	s.auditRepo.Log(ctx, &models.AuditLog{
		CompanyID: user.CompanyID,
		UserID:    &user.ID,
		UserEmail: &user.Email,
		Action:    "USER_LOGIN",
		Resource:  "auth",
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		Severity:  "INFO",
	})

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    s.cfg.AccessExpiryHours * 3600,
		User:         user,
	}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, rawRefreshToken, ipAddress, userAgent string) (*LoginResponse, error) {
	tokenHash := hashToken(rawRefreshToken)

	tokenRecord, err := s.tokenRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if tokenRecord.RevokedAt != nil {
		return nil, ErrTokenRevoked
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, ErrTokenInvalid
	}

	// Revoke old token (rotation)
	s.tokenRepo.RevokeRefreshToken(ctx, tokenHash)

	user, err := s.userRepo.FindByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	roles, _ := s.userRepo.GetUserRoles(ctx, user.ID)
	user.Roles = roles

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	rawRefresh, newHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(s.cfg.RefreshExpiryDays) * 24 * time.Hour)
	s.tokenRepo.StoreRefreshToken(ctx, user.ID, newHash, expiresAt, ipAddress, userAgent)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    s.cfg.AccessExpiryHours * 3600,
		User:         user,
	}, nil
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string, userID uuid.UUID, ip ...string) error {
	if rawRefreshToken != "" {
		tokenHash := hashToken(rawRefreshToken)
		s.tokenRepo.RevokeRefreshToken(ctx, tokenHash)
	}
	return nil
}

// RefreshToken is an alias for RefreshTokens to support both calling conventions
func (s *Service) RefreshToken(ctx context.Context, rawRefreshToken, ipAddress string) (*LoginResponse, error) {
	return s.RefreshTokens(ctx, rawRefreshToken, ipAddress, "")
}

// GetUserByID returns a user by their ID
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	roles, _ := s.userRepo.GetUserRoles(ctx, userID)
	user.Roles = roles
	return user, nil
}

func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.AccessSecret), nil
	})

	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenType != "access" {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (s *Service) generateAccessToken(user *models.User) (string, error) {
	roleNames := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleNames[i] = r.Code
	}

	claims := Claims{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Email:     user.Email,
		Roles:     roleNames,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.AccessExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "finsight",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.AccessSecret))
}

func generateRefreshToken() (raw, hash string, err error) {
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil {
		return
	}
	raw = hex.EncodeToString(bytes)
	hash = hashToken(raw)
	return
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
