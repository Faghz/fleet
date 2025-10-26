package api

import (
	"github.com/elzestia/fleet/configs"
	"github.com/gofiber/fiber/v2"

	inthttp "github.com/elzestia/fleet/pkg/transport/http"
	httphndl "github.com/elzestia/fleet/pkg/transport/http/handler"
	"github.com/elzestia/fleet/service"
	"go.uber.org/zap"
)

func CreateApiServer(cfg *configs.Config, logger *zap.Logger, services *service.Services) (httpServer *fiber.App) {
	httpServer = inthttp.CreateHttpServer(cfg.Http.Port, cfg.App.ContextTimeout, logger, cfg)
	httphndl.CreateHandler(httpServer, services, logger)
	inthttp.SetupAndServe(httpServer, services, cfg.Http.Port, cfg.App.ContextTimeout, logger)

	return
}

func ShutdownServer(httpServer *fiber.App, logger *zap.Logger) error {
	logger.Info("Shutting down server...")
	if err := httpServer.Shutdown(); err != nil {
		return err
	}

	logger.Info("Http Server shutdown successfully")

	return nil
}
