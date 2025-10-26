package httphndl

import (
	"net/http"
	"strings"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *HttpHandler) authMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return response.ResponseJson(c, http.StatusUnauthorized, "Missing Authorization Header")
		}

		ctx := c.UserContext()
		claims, err := h.services.UserService.VerifyToken(ctx, strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			h.logger.Info("Failed to verify token", zap.Error(err))
			return response.ResponseJson(c, http.StatusUnauthorized, "Invalid Token")
		}

		_, err = h.services.UserService.GetSessionByID(ctx, &claims)
		if err != nil {
			h.logger.Info("Failed to get session", zap.Error(err))
			return response.ResponseJson(c, http.StatusUnauthorized, "Invalid Session")
		}

		c.Locals("userID", claims.Subject)
		c.Locals("sessionID", claims.ID)
		c.Locals("orgID", claims.OrgID)

		return c.Next()
	}
}

func getUserClaims(c *fiber.Ctx) (claims *models.AuthClaims) {
	userID := c.Locals("userID")
	sessionID := c.Locals("sessionID")
	orgID := c.Locals("orgID")

	if userID == nil || sessionID == nil || orgID == nil {
		return
	}

	return &models.AuthClaims{
		Subject: userID.(string),
		ID:      sessionID.(string),
		OrgID:   orgID.(string),
	}
}
