package httphndl

import (
	"net/http"

	inthttp "github.com/elzestia/fleet/pkg/transport/http"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func createAuthHandler(app *fiber.App, handler *HttpHandler) {
	v1 := app.Group("/v1/auth")
	// Define routes
	v1.Post("/login", handler.login)
	v1.Post("/register", handler.registerUser)
	v1.Delete("/logout", handler.authMiddleware(), handler.logout)
}

// @Summary Login endpoint
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.Login true "Login request"
// @Success 200 {object} response.Login "Login success"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/auth/login [post]
func (h *HttpHandler) login(c *fiber.Ctx) error {
	req := request.Login{}
	if err := c.BodyParser(&req); err != nil {
		h.logger.Info("Failed to bind request", zap.Error(err))
		return response.ResponseJson(c, http.StatusNotAcceptable, "Invalid Request data")
	}

	validator := inthttp.GetValidator()
	if err := validator.Validate(&req); err != nil {
		return err
	}

	res, err := h.services.UserService.Login(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, http.StatusOK, "Login successful", res)
}

// @Summary Register user endpoint
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterUserRequest true "Register user request"
// @Success 200 {object} string "Register success"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/auth/register [post]
func (h *HttpHandler) registerUser(c *fiber.Ctx) error {
	req := request.RegisterUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		h.logger.Info("Failed to bind request", zap.Error(err))
		return err
	}

	validator := inthttp.GetValidator()
	if err := validator.Validate(&req); err != nil {
		return err
	}

	err := h.services.UserService.RegisterUser(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, http.StatusOK, "User Created")
}

// @Summary Logout endpoint
// @Description Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} string "Logout success"
// @Failure 400 {object} response.Failure "Invalid Request data"
// @Failure 500 {object} response.Failure "Internal Server Error"
// @Router /v1/auth/logout [delete]
func (h *HttpHandler) logout(c *fiber.Ctx) error {
	userClaims := getUserClaims(c)
	if userClaims == nil {
		return response.ResponseJson(c, http.StatusUnauthorized, "Unauthorized")
	}

	err := h.services.UserService.Logout(c.UserContext(), userClaims)
	if err != nil {
		return err
	}

	return response.ResponseJson(c, http.StatusOK, "Logout successful")
}
