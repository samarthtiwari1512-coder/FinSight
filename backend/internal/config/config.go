package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	FX       FXConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         string
	Mode         string // debug, release
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AllowOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpiryHours int
	RefreshExpiryDays int
}

type FXConfig struct {
	Provider      string // EXCHANGERATE_API, FIXER, OPEN_EXCHANGE
	APIKey        string
	UpdateIntervalMinutes int
	StalenessMinutes      int
}

type AppConfig struct {
	Environment  string
	LogLevel     string
	BaseCurrency string
	MaxPageSize  int
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			AllowOrigins: []string{
				getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
				"http://localhost:3000",
			},
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "finsight"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "finsight"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 25),
			MinConns: getEnvInt("DB_MIN_CONNS", 5),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:      requireEnv("JWT_ACCESS_SECRET"),
			RefreshSecret:     requireEnv("JWT_REFRESH_SECRET"),
			AccessExpiryHours: getEnvInt("JWT_ACCESS_EXPIRY_HOURS", 8),
			RefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 30),
		},
		FX: FXConfig{
			Provider:              getEnv("FX_PROVIDER", "EXCHANGERATE_API"),
			APIKey:                getEnv("FX_API_KEY", ""),
			UpdateIntervalMinutes: getEnvInt("FX_UPDATE_INTERVAL_MINUTES", 60),
			StalenessMinutes:      getEnvInt("FX_STALENESS_MINUTES", 120),
		},
		App: AppConfig{
			Environment:  getEnv("APP_ENV", "development"),
			LogLevel:     getEnv("LOG_LEVEL", "info"),
			BaseCurrency: getEnv("BASE_CURRENCY", "INR"),
			MaxPageSize:  getEnvInt("MAX_PAGE_SIZE", 200),
		},
	}

	return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	// Prefer DATABASE_URL if set (Neon, Render, Railway, etc.)
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d pool_min_conns=%d",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.MaxConns, c.MinConns,
	)
}

func (c *DatabaseConfig) LibPQDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// Config accessor methods required by FXService and workers

func (c *Config) GetFXAPIKey() string             { return c.FX.APIKey }
func (c *Config) GetFXProvider() string           { return c.FX.Provider }
func (c *Config) GetFXStalenessMinutes() int      { return c.FX.StalenessMinutes }
func (c *Config) GetBaseCurrency() string         { return c.App.BaseCurrency }
func (c *Config) GetFXUpdateIntervalMinutes() int { return c.FX.UpdateIntervalMinutes }

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		// In development, use defaults
		defaults := map[string]string{
			"JWT_ACCESS_SECRET":  "finsight-dev-access-secret-change-in-production",
			"JWT_REFRESH_SECRET": "finsight-dev-refresh-secret-change-in-production",
		}
		if d, ok := defaults[key]; ok {
			return d
		}
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return val
}
