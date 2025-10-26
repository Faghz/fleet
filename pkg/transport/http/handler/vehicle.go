package httphndl

import (
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func createVehicleHandler(app *fiber.App, handler *HttpHandler) {
	v1 := app.Group("/v1/vehicles")

	v1.Get("/:vehicle_id/location", handler.authMiddleware(), handler.getVehicleLocation)
}

// @Summary Get Vehicle Location
// @Description Get the current location of a vehicle by its ID
// @Tags vehicles
// @Produce json
// @Param vehicle_id path string true "Vehicle Entity ID"
// @Success 200 {object} response.VehicleLocation "Vehicle location retrieved successfully"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 404 {object} response.Failure "Vehicle not found"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/vehicles/{vehicle_id}/location [get]
func (h *HttpHandler) getVehicleLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	ctx := c.UserContext()

	vehicleLocation, err := h.services.VehicleService.GetVehicleLatestLocationByVehicleID(ctx, vehicleID)
	if err != nil {
		h.logger.Error("Failed to get vehicle location", zap.String("vehicle_id", vehicleID), zap.Error(err))
		return response.GenerateFailure(fiber.StatusInternalServerError, "Failed to get vehicle location", err.Error())
	}

	return response.ResponseJson(c, fiber.StatusOK, "Vehicle location retrieved successfully", vehicleLocation)
}
