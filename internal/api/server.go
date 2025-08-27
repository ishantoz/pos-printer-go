package api

import (
	"pos-printer/internal/config"
	"pos-printer/internal/db"
	printers "pos-printer/internal/printer"

	// "pos-printer/internal/job"
	// "pos-printer/internal/store"

	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo       *echo.Echo
	cfg        *config.Config
	sqlite     *db.SQLite
	posPrinter *printers.PosPrinter
}

func NewServer(cfg *config.Config, sqlite *db.SQLite, posPrinter *printers.PosPrinter) *Server {
	e := echo.New()

	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	srv := &Server{echo: e, cfg: cfg, sqlite: sqlite, posPrinter: posPrinter}
	srv.registerRoutes()
	return srv
}

func (server *Server) StartTLS() error {
	return server.echo.StartTLS(
		server.cfg.ServerConfig.Endpoint,
		server.cfg.ServerConfig.CertPath,
		server.cfg.ServerConfig.KeyPath,
	)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (server *Server) registerRoutes() {
	server.echo.GET("/health", server.healthCheckHandler)
	server.echo.POST("/barcode/print", server.printBarcodeHandler)
	server.echo.GET("/barcode/job/:id", server.jobBarcodeHandler)
}
