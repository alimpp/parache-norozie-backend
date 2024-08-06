package api

import (
	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func login(ctx *fiber.Ctx) error {

	return nil
}

func otp(ctx *fiber.Ctx) error {
	return nil
}

func password(ctx *fiber.Ctx) error {
	return nil
}
