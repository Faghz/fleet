package httphndl

import (
	"net/http"

	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/gofiber/fiber/v2"
)

func createUserHandler(app *fiber.App, handler *HttpHandler) {
	v1 := app.Group("/v1/users")
	// Define routes
	v1.Get("/me", handler.authMiddleware(), handler.getUserDetail)
}

// @Summary Get current user details
// @Description Get details of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.UserDetail "User details retrieved successfully"
// @Failure 401 {object} response.Failure "Unauthorized"
// @Failure 404 {object} response.Failure "User not found"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/users/me [get]
func (h *HttpHandler) getUserDetail(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return response.ResponseJson(c, http.StatusUnauthorized, "Unauthorized")
	}

	userDetail, err := h.services.UserService.GetUserByUserID(c.UserContext(), claims.Subject)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, http.StatusOK, "User details retrieved successfully", userDetail)
}
