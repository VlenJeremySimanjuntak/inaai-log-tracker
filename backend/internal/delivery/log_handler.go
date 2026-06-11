package delivery

import (
	"net/http"
	"strconv"
	"time"
	"context"

	"backend-tracker/internal/models"
	"backend-tracker/internal/usecase"

	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	usecase usecase.LogUsecase
}

func NewLogHandler(e *echo.Echo, u usecase.LogUsecase) {
	h := &LogHandler{usecase: u}
	g := e.Group("/api")
	g.GET("/logs", h.GetAllLogs)
	g.POST("/logs", h.CreateLog)
	g.PUT("/logs/:id/status", h.UpdateStatus)
	g.GET("/ai-summary", h.GetSummary)
}

func (h *LogHandler) GetAllLogs(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	logs, err := h.usecase.GetAllLogs(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, logs)
}

func (h *LogHandler) CreateLog(c echo.Context) error {
	var input models.IncidentLog
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
	}
	if input.UserID == 0 || input.CategoryID == 0 || input.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id, category_id, title wajib diisi"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	err := h.usecase.ReportIncident(ctx, &input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, input)
}

func (h *LogHandler) UpdateStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
	}
	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload status tidak valid"})
	}
	newStatus, ok := body["status"]
	if !ok || newStatus == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Field status wajib diisi"})
	}
	allowed := map[string]bool{"Menunggu": true, "Diproses": true, "Selesai": true}
	if !allowed[newStatus] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status tidak valid (Menunggu, Diproses, Selesai)"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()
	err = h.usecase.ChangeLogStatus(ctx, id, newStatus)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Log tidak ditemukan"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Status berhasil diubah dengan pessimistic locking"})
}

func (h *LogHandler) GetSummary(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()
	summary, err := h.usecase.GetOrTriggerAggregateSummary(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}