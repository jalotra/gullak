package main

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mr-karan/gullak/internal/db"
	"github.com/mr-karan/gullak/internal/llm"
)

// TODO: Fix the API response to return the correct status codes and messages with content types.

type App struct {
	srv     *echo.Echo
	log     *slog.Logger
	addr    string
	llm     *llm.Manager
	queries *db.Queries
}

func initApp(addr string, timeout time.Duration, static fs.FS, queries *db.Queries, llmMgr *llm.Manager, log *slog.Logger) *App {
	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.Logger()) -> Too noisy for now.
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: timeout,
	}))

	// Register handlers.

	e.GET("/api", handleIndex)                                               // Simple welcome message or API status
	e.POST("/api/transactions", handleCreateTransaction)                     // Creates a new transaction
	e.GET("/api/transactions", handleListTransactions)                       // Lists all transactions, with optional filters
	e.GET("/api/transactions/:id", handleGetTransaction)                     // Retrieves a specific transaction by ID
	e.PUT("/api/transactions/:id", handleUpdateTransaction)                  // Updates a specific transaction by ID
	e.DELETE("/api/transactions/:id", handleDeleteTransaction)               // Deletes a specific transaction by ID
	e.GET("/api/reports/top-expense-categories", handleTopExpenseCategories) // Retrieves top expense categories
	e.GET("/api/reports/daily-spending", handleDailySpending)                // Retrieves spending for a specific day
	// e.GET("/api/reports/monthly-spending-summary", handleMonthlySpendingSummary) // Retrieves spending summary by month

	// Middleware to serve the static files.
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/",
		Index:      "index.html",
		HTML5:      true, // This kicks in client side routing.
		Filesystem: http.FS(static),
	}))

	return &App{
		srv:     e,
		log:     log,
		addr:    addr,
		queries: queries,
		llm:     llmMgr,
	}
}

func (m *App) Start(ctx context.Context) error {
	// Register app (*App) to be injected into all HTTP handlers.
	m.srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", m)
			return next(c)
		}
	})

	// Start server in a goroutine to allow for graceful shutdown.
	go func() {
		if err := m.srv.Start(m.addr); err != http.ErrServerClosed {
			m.srv.Logger.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := m.srv.Shutdown(shutdownCtx); err != nil {
		m.srv.Logger.Fatalf("Error shutting down server: %v", err)
	}

	return nil
}
