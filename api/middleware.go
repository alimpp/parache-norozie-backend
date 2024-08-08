package api

import (
	"ecom/pkg/constants"
	"errors"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"time"
)

func rateLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        240,             // Maximum number of requests
		Expiration: 1 * time.Minute, // Time window in seconds
	})
}

func logMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "{\"level\":\"info\",\"service\":\"parche-norozie-backend\",\"error\":\"${status}\",\"time\":\"${time}\",\"message\":\"${locals:requestid} ${latency} ${method} ${path}\"}\n",
		TimeFormat: "2006-01-02T15:04:05-0700",
	})
}

func requestIdMiddleware() fiber.Handler {
	return requestid.New()
}

func performanceMonitorMiddleware() fiber.Handler {
	return pprof.New()
}

func PrometheusMiddleware(app *fiber.App) fiber.Handler {
	prometheus := fiberprometheus.New(constants.ServiceName)
	prometheus.RegisterAt(app, "/metrics")
	return prometheus.Middleware
}

func CacheHeaderMiddleware(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public, max-age=21600")
	return c.Next()
}

func ErrorResponse(ctx *fiber.Ctx, err error) error {
	constants.Logger.Err(err).Msgf(err.Error())

	errorResp := Resp{Message: err.Error()}

	switch {
	case errors.Is(err, fiber.ErrUnprocessableEntity):
		errorResp.Message = "invalid body request"
		errorResp.Status = fiber.StatusUnprocessableEntity

	case errors.Is(err, constants.ErrRequestBody{}):
		errorResp.Message = "invalid body request"
		errorResp.Status = fiber.StatusBadRequest
	default:
		errorResp.Message = "unknown error"
		errorResp.Status = fiber.StatusInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(errorResp)
}
