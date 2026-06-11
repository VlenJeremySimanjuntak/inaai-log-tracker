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

// LogHandler adalah handler untuk endpoint log gangguan
type LogHandler struct {
	usecase usecase.LogUsecase
}

// NewLogHandler mendaftarkan semua route ke Echo
func NewLogHandler(e *echo.Echo, u usecase.LogUsecase) {
	h := &LogHandler{usecase: u}
	g := e.Group("/api")
	g.GET("/logs", h.GetAllLogs)
	g.POST("/logs", h.CreateLog)
	g.PUT("/logs/:id/status", h.UpdateStatus)
	g.GET("/ai-summary", h.GetSummary)
}

// GetAllLogs godoc
// @Summary      Dapatkan semua laporan gangguan
// @Description  Mengembalikan daftar semua incident log yang sudah tersimpan, diurutkan dari yang terbaru
// @Tags         Logs
// @Produce      json
// @Success      200  {array}   models.IncidentLog
// @Failure      500  {object}  map[string]string
// @Router       /api/logs [get]
func (h *LogHandler) GetAllLogs(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	logs, err := h.usecase.GetAllLogs(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, logs)
}

// CreateLog godoc
// @Summary      Buat laporan gangguan baru
// @Description  Menambahkan incident log baru dengan status awal "Menunggu"
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param        log  body  models.IncidentLog  true  "Data laporan (user_id, category_id, title, description)"
// @Success      201  {object}  models.IncidentLog
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/logs [post]
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

// UpdateStatus godoc
// @Summary      Ubah status laporan gangguan
// @Description  Mengubah status incident log (Menunggu, Diproses, Selesai) menggunakan pessimistic locking untuk mencegah race condition.
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param        id     path      int    true  "ID laporan"
// @Param        status body      string true  "Status baru (Menunggu/Diproses/Selesai)"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/logs/{id}/status [put]
func (h *LogHandler) UpdateStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
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
		// Deteksi not found berdasarkan pesan error (bisa diperbaiki dengan error sentinel)
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Log tidak ditemukan"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Status berhasil diubah dengan pessimistic locking"})
}

// GetSummary godoc
// @Summary      Dapatkan ringkasan AI
// @Description  Mengembalikan ringkasan otomatis dari semua laporan menggunakan Gemini AI. Jika belum ada atau ada laporan baru, akan memicu pembuatan ringkasan baru.
// @Tags         AI
// @Produce      json
// @Success      200  {object}  models.AISummary
// @Failure      500  {object}  map[string]string
// @Router       /api/ai-summary [get]
func (h *LogHandler) GetSummary(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()
	summary, err := h.usecase.GetOrTriggerAggregateSummary(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}