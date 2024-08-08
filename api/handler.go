package api

import (
	_ "ecom/cmd/docs"
	"ecom/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// @Summary      User otp login
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        login body LoginRes true "Login credentials"
// @Success      200 {object} Resp
// @Router       /login [post]
func login(ctx *fiber.Ctx) error {
	req := LoginReq{}
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, err)
	}

	if err := req.Validate(); err != nil {
		return ErrorResponse(ctx, err)
	}

	otpCode := util.GenerateRandomString(5)
	err := AppSrv.sms.Send(req.Phone, map[string]string{"code: ": otpCode})
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	resp := LoginResp{UserExists: true, TTE: AppSrv.cfg.OTP.TTE}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func verifyOtp(ctx *fiber.Ctx) error {
	return nil
}

func password(ctx *fiber.Ctx) error {
	return nil
}
