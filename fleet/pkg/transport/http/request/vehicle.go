package request

import "github.com/gofiber/fiber/v2"

type GetVehicleLocationHistoryRequest struct {
	VehicleID string `validate:"required"`
	Start     int64  `validate:"required"`
	End       int64  `validate:"required"`
}

func (r *GetVehicleLocationHistoryRequest) ParseFromContext(c *fiber.Ctx) {
	r.VehicleID = c.Params("vehicle_id")
	r.Start = int64(c.QueryInt("start", 0))
	r.End = int64(c.QueryInt("end", 0))
}
