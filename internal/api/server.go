package api

import (
	"pos-printer/internal/config"

	"pos-printer/internal/api/handlers"

	// "pos-printer/internal/job"
	// "pos-printer/internal/store"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	// store     *store.Store
	// processor *job.Processor
}

func NewServer() *Server {
	e := echo.New()

	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	srv := &Server{echo: e}
	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {

	// Health Check
	s.echo.GET("/health", handlers.HealthCheckHandler)

	// Barcode Print API
	s.echo.POST("/barcode/print", handlers.EnqueueBarcodeHandler)
	// s.echo.GET("/barcode/job-status/:id", jobBarcodeStatusHandler)
}

func (s *Server) StartTLS(cfg *config.ServerConfig) error {
	return s.echo.StartTLS(cfg.Endpoint, cfg.CertPath, cfg.KeyPath)
}
