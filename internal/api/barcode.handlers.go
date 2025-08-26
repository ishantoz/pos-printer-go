package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"pos-printer/internal/model"

	"github.com/labstack/echo/v4"
)

func (server *Server) healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) printBarcodeHandler(c echo.Context) error {
	var req model.PrintBarcodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			echo.Map{
				"error": "Invalid JSON",
			},
		)
	}

	server.applyDefaultsBarcodeHelper(&req)
	if err := server.validateBarcodeRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			echo.Map{
				"error": err.Error(),
			},
		)
	}

	if err := server.posPrinter.CheckPrinter(req.VID, req.PID); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			echo.Map{
				"error": fmt.Sprintf(
					"Printer device not found, please check connected or not: %s",
					err,
				),
			},
		)
	}

	jobId, err := server.sqlite.EnqueueBarcodeJob(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to enqueue job"})
	}

	return c.JSON(http.StatusAccepted,
		echo.Map{
			"jobId":  jobId,
			"status": server.cfg.WorkerConfig.JobStatus.StatusPending,
		},
	)
}

func (server *Server) jobBarcodeHandler(c echo.Context) error {
	id := c.Param("id")

	job, err := server.sqlite.FetchBarcodeJob(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error fetching job"})
	}

	return c.JSON(http.StatusOK, job)
}
