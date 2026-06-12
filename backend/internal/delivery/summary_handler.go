// internal/delivery/summary_handler.go
package delivery

import (
	"net/http"
	"backend-tracker/internal/usecase"
	"github.com/labstack/echo/v4"
)

type SummaryHandler struct {
	uc usecase.LogUsecase
}

func NewSummaryHandler(uc usecase.LogUsecase) *SummaryHandler {
	return &SummaryHandler{uc: uc}
}

func (h *SummaryHandler) GetLatestSummary(c echo.Context) error {
	summary, err := h.uc.GetOrTriggerAggregateSummary(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}