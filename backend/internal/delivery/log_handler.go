package delivery

import (
	"net/http"
	"strconv"
	"backend-tracker/internal/models"
	"backend-tracker/internal/usecase"
	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	Usecase usecase.LogUsecase
}

func NewLogHandler(e *echo.Echo, u usecase.LogUsecase) {
	handler := &LogHandler{Usecase: u}

	e.GET("/api/logs", handler.GetAllLogs)
	e.POST("/api/logs", handler.CreateLog)
	e.PUT("/api/logs/:id/status", handler.UpdateStatus)
	e.GET("/api/ai-summary", handler.GetSummary)
}

func (h *LogHandler) GetAllLogs(c echo.Context) error {
	logs, err := h.Usecase.GetAllLogs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, logs)
}

func (h *LogHandler) CreateLog(c echo.Context) error {
	var input models.IncidentLog
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload data tidak valid"})
	}

	if err := h.Usecase.ReportIncident(c.Request().Context(), &input); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, input)
}

func (h *LogHandler) UpdateStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID Laporan salah"})
	}

	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload status tidak valid"})
	}

	newStatus := body["status"]
	if err := h.Usecase.ChangeLogStatus(c.Request().Context(), id, newStatus); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Status laporan berhasil diubah melalui antrean aman"})
}

func (h *LogHandler) GetSummary(c echo.Context) error {
	summary, err := h.Usecase.GetOrTriggerAggregateSummary(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}