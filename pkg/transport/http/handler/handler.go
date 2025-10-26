package httphndl

import (
	"net/http"

	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/elzestia/fleet/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HttpHandler struct {
	services *service.Services
	logger   *zap.Logger
}

func CreateHandler(app *fiber.App, services *service.Services, logger *zap.Logger) {
	handler := &HttpHandler{
		services: services,
		logger:   logger,
	}

	// Define routes
	app.Get("/service/healthz", handler.healthz)
	createAuthHandler(app, handler)
	createUserHandler(app, handler)
}

// Service Health check
// @Summary Service Health check
// @Description Service Health check
// @Tags service
// @Produce json
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.Failure
// @Failure 409 {object} response.Failure
// @Failure 500 {object} response.Failure
// @Router /service/healthz [get]
func (h *HttpHandler) healthz(c *fiber.Ctx) error {
	ctx := c.UserContext()

	if err := h.services.HealthzService.Healthz(ctx); err != nil {
		h.logger.Error("Health check failed", zap.Error(err))
		return response.GenerateFailure(http.StatusServiceUnavailable, "Service Unavailable", "Health check failed: "+err.Error())
	}

	return response.ResponseJson(c, http.StatusOK, "Service is up and running")
}
