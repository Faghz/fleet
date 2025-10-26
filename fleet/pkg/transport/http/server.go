package inthttp

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"

	"go.uber.org/zap"

	_ "github.com/elzestia/fleet/cmd/server/docs"
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/service"
)

func CreateHttpServer(port string, timeout time.Duration, logger *zap.Logger, config *configs.Config) *fiber.App {
	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		ErrorHandler: customErrorHandler,
	})

	// Initialize custom validator
	CreateCustomValidator("EN")

	setupMiddleware(fiberApp, logger)
	setupSwagger(fiberApp)
	setupCORS(fiberApp, config.Http.AllowedOrigins, config.Http.AllowCredentials)

	return fiberApp
}

func SetupAndServe(
	fiberApp *fiber.App,
	services *service.Services,
	port string,
	timeout time.Duration,
	logger *zap.Logger,
) {
	if err := fiberApp.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

}

func ShutdownServer(fiberApp *fiber.App) error {
	if err := fiberApp.Shutdown(); err != nil {
		return err
	}

	return nil
}

func setupMiddleware(app *fiber.App, logger *zap.Logger) {
	// Panic recovery middleware
	app.Use(recover.New())

	// Request ID middleware
	app.Use(requestid.New())

	// Custom logging middleware
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request
		latency := time.Since(start)
		logAttributes := []zap.Field{
			zap.String("latency", latency.String()),
			zap.String("requestId", c.GetRespHeader("X-Request-ID")),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.String("uri", c.OriginalURL()),
			zap.Int("status", c.Response().StatusCode()),
		}

		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			logger.Error("error", logAttributes...)
		case status >= 400:
			logger.Warn("bad-request", logAttributes...)
		case status >= 300:
			logger.Warn("success", logAttributes...)
		default:
			logger.Info("success", logAttributes...)
		}

		return err
	})
}

func setupSwagger(fiberApp *fiber.App) {
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)
}

func setupCORS(app *fiber.App, allowedOrigins string, allowCredentials bool) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Authorization,Access-Control-Allow-Origin,token,Pv,Content-Type,Accept,Content-Length,Accept-Encoding,X-CSRF-Token",
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin",
		AllowCredentials: allowCredentials,
	}))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": message,
		"data":    nil,
	})
}
