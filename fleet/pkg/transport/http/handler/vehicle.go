package httphndl

import (
	inthttp "github.com/elzestia/fleet/pkg/transport/http"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/gofiber/fiber/v2"
)

func createVehicleHandler(app *fiber.App, handler *HttpHandler) {
	v1 := app.Group("/v1/vehicles")

	v1.Get("/:vehicle_id/location", handler.authMiddleware(), handler.getVehicleLocation)
	v1.Get("/:vehicle_id/history", handler.authMiddleware(), handler.getVehicleLocationHistory)
}

// @Summary Get Vehicle Location
// @Description Get the current location of a vehicle by its ID
// @Tags vehicles
// @Produce json
// @Param vehicle_id path string true "Vehicle Entity ID"
// @Success 200 {object} response.BaseResponse{data=response.VehicleLocation} "Vehicle location retrieved successfully"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 404 {object} response.Failure "Vehicle not found"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/vehicles/{vehicle_id}/location [get]
func (h *HttpHandler) getVehicleLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	ctx := c.UserContext()

	vehicleLocation, err := h.services.VehicleService.GetVehicleLatestLocationByVehicleID(ctx, vehicleID)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, fiber.StatusOK, "Vehicle location retrieved successfully", vehicleLocation)
}

// @Summary Get Vehicle Location History
// @Description Get the location history of a vehicle by its ID
// @Tags vehicles
// @Produce json
// @Param vehicle_id path string true "Vehicle Entity ID"
// @Param start_time query int64 false "Start time for location history filter (Unix timestamp)"
// @Param end_time query int64 false "End time for location history filter (Unix timestamp)"
// @Success 200 {array} response.BaseResponse{data=response.VehicleLocation} "Vehicle location history retrieved successfully"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 404 {object} response.Failure "Vehicle not found"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/vehicles/{vehicle_id}/locations/history [get]
func (h *HttpHandler) getVehicleLocationHistory(c *fiber.Ctx) error {
	req := request.GetVehicleLocationHistoryRequest{}
	req.ParseFromContext(c)

	err := inthttp.GetValidator().Validate(&req)
	if err != nil {
		return response.FailureResponse(c, err)
	}

	vehicleLocations, err := h.services.VehicleService.GetVehicleLocationHistory(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, fiber.StatusOK, "Vehicle location history retrieved successfully", vehicleLocations)
}
