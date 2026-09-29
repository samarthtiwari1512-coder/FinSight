package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/finsight/backend/internal/auth"
	"github.com/finsight/backend/internal/config"
	"github.com/finsight/backend/internal/handlers"
	"github.com/finsight/backend/internal/middleware"
	"github.com/finsight/backend/internal/repositories"
	"github.com/finsight/backend/internal/services"
	"github.com/finsight/backend/internal/workers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load .env for local development
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Logger
	var logger *zap.Logger
	if cfg.App.Environment == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Database
	dbPool, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		logger.Fatal("Database ping failed", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Redis ping failed, some features will be degraded", zap.Error(err))
	} else {
		logger.Info("Connected to Redis")
	}

	// Repositories
	userRepo := repositories.NewUserRepository(dbPool)
	tokenRepo := repositories.NewTokenRepository(dbPool)
	auditRepo := repositories.NewAuditRepository(dbPool)
	companyRepo := repositories.NewCompanyRepository(dbPool)
	entityRepo := repositories.NewEntityRepository(dbPool)
	customerRepo := repositories.NewCustomerRepository(dbPool)
	supplierRepo := repositories.NewSupplierRepository(dbPool)
	invoiceRepo := repositories.NewInvoiceRepository(dbPool)
	billRepo := repositories.NewBillRepository(dbPool)
	paymentRepo := repositories.NewPaymentRepository(dbPool)
	bankRepo := repositories.NewBankRepository(dbPool)
	fxRepo := repositories.NewFXRepository(dbPool)
	alertRepo := repositories.NewAlertRepository(dbPool)
	scenarioRepo := repositories.NewScenarioRepository(dbPool)
	riskRepo := repositories.NewRiskRepository(dbPool)
	reconciliationRepo := repositories.NewReconciliationRepository(dbPool)
	cashRepo := repositories.NewCashRepository(dbPool)

	// Services
	authSvc := auth.NewService(userRepo, tokenRepo, auditRepo, &cfg.JWT)
	fxSvc := services.NewFXService(fxRepo, rdb, cfg, logger)
	wcSvc := services.NewWorkingCapitalService(invoiceRepo, billRepo, fxSvc, cfg)
	cashSvc := services.NewCashService(cashRepo, bankRepo, fxSvc, cfg)
	forecastSvc := services.NewForecastService(invoiceRepo, billRepo, cashRepo, fxSvc, cfg, logger)
	alertSvc := services.NewAlertService(alertRepo, rdb, logger)
	dashSvc := services.NewDashboardService(cashSvc, wcSvc, alertSvc, fxSvc, invoiceRepo, billRepo, cfg)
	riskSvc := services.NewRiskService(riskRepo, fxSvc, alertSvc, cfg, logger)
	scenarioSvc := services.NewScenarioService(scenarioRepo, cashSvc, wcSvc, fxSvc, cfg)
	reconSvc := services.NewReconciliationService(reconciliationRepo, bankRepo, auditRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc, auditRepo)
	dashHandler := handlers.NewDashboardHandler(dashSvc)
	cashHandler := handlers.NewCashHandler(cashSvc, bankRepo, fxSvc)
	customerHandler := handlers.NewCustomerHandler(customerRepo, invoiceRepo, auditRepo)
	supplierHandler := handlers.NewSupplierHandler(supplierRepo, billRepo, auditRepo)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceRepo, customerRepo, fxSvc, auditRepo, wcSvc)
	billHandler := handlers.NewBillHandler(billRepo, supplierRepo, fxSvc, auditRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentRepo, fxSvc, auditRepo, cfg)
	bankHandler := handlers.NewBankHandler(bankRepo, auditRepo)
	fxHandler := handlers.NewFXHandler(fxSvc, fxRepo, auditRepo)
	wcHandler := handlers.NewWorkingCapitalHandler(wcSvc)
	forecastHandler := handlers.NewForecastHandler(forecastSvc)
	alertHandler := handlers.NewAlertHandler(alertSvc)
	scenarioHandler := handlers.NewScenarioHandler(scenarioSvc, scenarioRepo)
	riskHandler := handlers.NewRiskHandler(riskSvc, riskRepo)
	reconHandler := handlers.NewReconciliationHandler(reconSvc, bankRepo)
	auditHandler := handlers.NewAuditHandler(auditRepo)
	reportHandler := handlers.NewReportHandler(invoiceRepo, billRepo, cashSvc, wcSvc, fxSvc)
	adminHandler := handlers.NewAdminHandler(userRepo, companyRepo, entityRepo, auditRepo)
	healthHandler := handlers.NewHealthHandler(dbPool, rdb, fxSvc)

	// Background workers
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	go workers.StartFXRateSyncer(workerCtx, fxSvc, cfg, logger)
	go workers.StartAlertEvaluator(workerCtx, alertSvc, invoiceRepo, billRepo, cashSvc, riskSvc, cfg, logger)
	go workers.StartForecastGenerator(workerCtx, forecastSvc, logger)

	// Gin setup
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.SecureHeaders())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Request-ID", "X-Idempotency-Key"},
		ExposeHeaders:    []string{"X-Request-ID", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health & metrics (no auth required)
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1
	v1 := r.Group("/api/v1")

	// Auth routes (no auth required)
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", middleware.Auth(authSvc), authHandler.Logout)
		authRoutes.GET("/me", middleware.Auth(authSvc), authHandler.Me)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.Auth(authSvc))
	{
		// Dashboard
		protected.GET("/dashboard", dashHandler.GetDashboard)
		protected.GET("/dashboard/summary", dashHandler.GetSummary)

		// Cash Management
		cash := protected.Group("/cash")
		{
			cash.GET("/position", cashHandler.GetPosition)
			cash.GET("/position/by-bank", cashHandler.GetPositionByBank)
			cash.GET("/position/by-currency", cashHandler.GetPositionByCurrency)
			cash.GET("/position/by-entity", cashHandler.GetPositionByEntity)
			cash.GET("/movements", cashHandler.GetMovements)
			cash.GET("/accounts", cashHandler.GetCashAccounts)
		}

		// Bank Accounts
		banks := protected.Group("/banks")
		{
			banks.GET("", bankHandler.ListBankAccounts)
			banks.GET("/:id", bankHandler.GetBankAccount)
			banks.POST("", middleware.RequirePermission("bank_accounts", "create"), bankHandler.CreateBankAccount)
			banks.PUT("/:id", middleware.RequirePermission("bank_accounts", "update"), bankHandler.UpdateBankAccount)
			banks.GET("/:id/transactions", bankHandler.GetTransactions)
			banks.POST("/:id/import", middleware.RequirePermission("bank_accounts", "update"), bankHandler.ImportStatement)
		}

		// Customers
		customers := protected.Group("/customers")
		{
			customers.GET("", customerHandler.ListCustomers)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.POST("", middleware.RequirePermission("customers", "create"), customerHandler.CreateCustomer)
			customers.PUT("/:id", middleware.RequirePermission("customers", "update"), customerHandler.UpdateCustomer)
			customers.GET("/:id/invoices", customerHandler.GetCustomerInvoices)
			customers.GET("/:id/risk", customerHandler.GetCustomerRisk)
			customers.GET("/concentration", customerHandler.GetConcentration)
		}

		// Suppliers
		suppliers := protected.Group("/suppliers")
		{
			suppliers.GET("", supplierHandler.ListSuppliers)
			suppliers.GET("/:id", supplierHandler.GetSupplier)
			suppliers.POST("", middleware.RequirePermission("suppliers", "create"), supplierHandler.CreateSupplier)
			suppliers.PUT("/:id", middleware.RequirePermission("suppliers", "update"), supplierHandler.UpdateSupplier)
			suppliers.GET("/:id/bills", supplierHandler.GetSupplierBills)
			suppliers.GET("/concentration", supplierHandler.GetConcentration)
		}

		// Invoices (AR)
		invoices := protected.Group("/invoices")
		{
			invoices.GET("", invoiceHandler.ListInvoices)
			invoices.GET("/:id", invoiceHandler.GetInvoice)
			invoices.POST("", middleware.RequirePermission("invoices", "create"), invoiceHandler.CreateInvoice)
			invoices.PUT("/:id", middleware.RequirePermission("invoices", "update"), invoiceHandler.UpdateInvoice)
			invoices.POST("/:id/approve", middleware.RequirePermission("invoices", "approve"), invoiceHandler.ApproveInvoice)
			invoices.POST("/:id/payment", middleware.RequirePermission("invoices", "update"), invoiceHandler.RecordPayment)
			invoices.POST("/:id/dispute", middleware.RequirePermission("invoices", "update"), invoiceHandler.DisputeInvoice)
			invoices.GET("/aging", invoiceHandler.GetAging)
			invoices.GET("/overdue", invoiceHandler.GetOverdue)
			invoices.GET("/export", middleware.RequirePermission("invoices", "export"), invoiceHandler.Export)
		}

		// Receivables summary
		receivables := protected.Group("/receivables")
		{
			receivables.GET("/summary", invoiceHandler.GetARSummary)
			receivables.GET("/aging", invoiceHandler.GetAging)
			receivables.GET("/collections", invoiceHandler.GetCollections)
		}

		// Bills (AP)
		bills := protected.Group("/bills")
		{
			bills.GET("", billHandler.ListBills)
			bills.GET("/:id", billHandler.GetBill)
			bills.POST("", middleware.RequirePermission("bills", "create"), billHandler.CreateBill)
			bills.PUT("/:id", middleware.RequirePermission("bills", "update"), billHandler.UpdateBill)
			bills.POST("/:id/approve", middleware.RequirePermission("bills", "approve"), billHandler.ApproveBill)
			bills.POST("/:id/payment", middleware.RequirePermission("bills", "update"), billHandler.RecordPayment)
			bills.GET("/aging", billHandler.GetAging)
			bills.GET("/overdue", billHandler.GetOverdue)
		}

		// Payables summary
		payables := protected.Group("/payables")
		{
			payables.GET("/summary", billHandler.GetAPSummary)
			payables.GET("/aging", billHandler.GetAging)
			payables.GET("/due-soon", billHandler.GetDueSoon)
		}

		// Working Capital
		wc := protected.Group("/working-capital")
		{
			wc.GET("/kpis", wcHandler.GetKPIs)
			wc.GET("/trend", wcHandler.GetTrend)
			wc.GET("/opportunities", wcHandler.GetOpportunities)
			wc.GET("/ccc", wcHandler.GetCCC)
		}

		// Cash Flow Forecast
		forecast := protected.Group("/forecast")
		{
			forecast.GET("", forecastHandler.GetForecast)
			forecast.GET("/7d", forecastHandler.GetForecast7d)
			forecast.GET("/30d", forecastHandler.GetForecast30d)
			forecast.GET("/90d", forecastHandler.GetForecast90d)
			forecast.GET("/accuracy", forecastHandler.GetAccuracy)
			forecast.POST("/generate", middleware.RequirePermission("forecast", "create"), forecastHandler.Generate)
		}

		// FX Management
		fx := protected.Group("/fx")
		{
			fx.GET("/rates", fxHandler.GetRates)
			fx.GET("/rates/:from/:to", fxHandler.GetRate)
			fx.GET("/rates/history", fxHandler.GetRateHistory)
			fx.GET("/exposure", fxHandler.GetExposure)
			fx.GET("/exposure/summary", fxHandler.GetExposureSummary)
			fx.POST("/sensitivity", fxHandler.RunSensitivity)
			fx.POST("/scenarios", fxHandler.RunScenario)

			// FX Deals
			deals := fx.Group("/deals")
			{
				deals.GET("", fxHandler.ListDeals)
				deals.GET("/:id", fxHandler.GetDeal)
				deals.POST("", middleware.RequirePermission("fx", "create"), fxHandler.CreateDeal)
				deals.POST("/:id/approve", middleware.RequirePermission("fx", "approve"), fxHandler.ApproveDeal)
			}

			// Hedges
			hedges := fx.Group("/hedges")
			{
				hedges.GET("", fxHandler.ListHedges)
				hedges.POST("", middleware.RequirePermission("fx", "create"), fxHandler.CreateHedge)
				hedges.GET("/coverage", fxHandler.GetHedgeCoverage)
			}
		}

		// Payments
		payments := protected.Group("/payments")
		{
			payments.GET("", paymentHandler.ListPayments)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.POST("", middleware.RequirePermission("payments", "create"), paymentHandler.CreatePayment)
			payments.POST("/:id/approve", middleware.RequirePermission("payments", "approve"), paymentHandler.ApprovePayment)
			payments.POST("/:id/reject", middleware.RequirePermission("payments", "approve"), paymentHandler.RejectPayment)
			payments.DELETE("/:id", paymentHandler.CancelPayment)
		}

		// Risk
		risk := protected.Group("/risk")
		{
			risk.GET("/events", riskHandler.ListEvents)
			risk.GET("/events/:id", riskHandler.GetEvent)
			risk.POST("/events/:id/acknowledge", riskHandler.AcknowledgeEvent)
			risk.POST("/events/:id/resolve", riskHandler.ResolveEvent)
			risk.GET("/limits", riskHandler.ListLimits)
			risk.POST("/limits", middleware.RequirePermission("risk", "manage"), riskHandler.CreateLimit)
			risk.PUT("/limits/:id", middleware.RequirePermission("risk", "manage"), riskHandler.UpdateLimit)
			risk.GET("/summary", riskHandler.GetRiskSummary)
		}

		// Scenarios
		scenarios := protected.Group("/scenarios")
		{
			scenarios.GET("", scenarioHandler.ListScenarios)
			scenarios.GET("/:id", scenarioHandler.GetScenario)
			scenarios.POST("", middleware.RequirePermission("scenarios", "create"), scenarioHandler.CreateScenario)
			scenarios.POST("/:id/run", scenarioHandler.RunScenario)
			scenarios.GET("/:id/results", scenarioHandler.GetResults)
			scenarios.POST("/stress-test", scenarioHandler.RunStressTest)
		}

		// Reconciliation
		recon := protected.Group("/reconciliation")
		{
			recon.GET("", reconHandler.ListReconciliations)
			recon.GET("/dashboard", reconHandler.GetDashboard)
			recon.POST("/match", middleware.RequirePermission("reconciliation", "create"), reconHandler.ManualMatch)
			recon.GET("/unmatched", reconHandler.GetUnmatched)
			recon.GET("/exceptions", reconHandler.GetExceptions)
		}

		// Alerts
		alerts := protected.Group("/alerts")
		{
			alerts.GET("", alertHandler.ListAlerts)
			alerts.GET("/:id", alertHandler.GetAlert)
			alerts.POST("/:id/acknowledge", alertHandler.AcknowledgeAlert)
			alerts.POST("/:id/dismiss", alertHandler.DismissAlert)
			alerts.GET("/counts", alertHandler.GetAlertCounts)
		}

		// Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", alertHandler.ListNotifications)
			notifications.POST("/:id/read", alertHandler.MarkRead)
			notifications.POST("/read-all", alertHandler.MarkAllRead)
			notifications.GET("/unread-count", alertHandler.GetUnreadCount)
		}

		// Reports
		reports := protected.Group("/reports")
		{
			reports.GET("", reportHandler.ListReports)
			reports.POST("/cash-position", middleware.RequirePermission("reports", "create"), reportHandler.GenerateCashPosition)
			reports.POST("/working-capital", middleware.RequirePermission("reports", "create"), reportHandler.GenerateWorkingCapital)
			reports.POST("/ar-aging", middleware.RequirePermission("reports", "create"), reportHandler.GenerateARAgeing)
			reports.POST("/ap-aging", middleware.RequirePermission("reports", "create"), reportHandler.GenerateAPAgeing)
			reports.POST("/fx-exposure", middleware.RequirePermission("reports", "create"), reportHandler.GenerateFXExposure)
			reports.POST("/executive-summary", middleware.RequirePermission("reports", "create"), reportHandler.GenerateExecutiveSummary)
			reports.GET("/:id/download", reportHandler.DownloadReport)
		}

		// Audit Logs
		protected.GET("/audit-logs",
			middleware.RequirePermission("audit_logs", "read"),
			auditHandler.ListAuditLogs)
		protected.GET("/audit-logs/:resource/:id",
			middleware.RequirePermission("audit_logs", "read"),
			auditHandler.GetResourceHistory)

		// Admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("ADMIN", "CFO"))
		{
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.POST("/users/:id/roles", adminHandler.AssignRole)
			admin.GET("/roles", adminHandler.ListRoles)
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)
			admin.GET("/entities", adminHandler.ListEntities)
			admin.POST("/entities", adminHandler.CreateEntity)
		}

		// System Health (detailed — admin only)
		protected.GET("/system/health",
			middleware.RequireRole("ADMIN"),
			healthHandler.DetailedHealth)
	}

	// Server setup
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("FinSight backend starting",
			zap.String("port", cfg.Server.Port),
			zap.String("env", cfg.App.Environment))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	workerCancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited cleanly")
}
