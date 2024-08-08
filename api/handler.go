package api

import (
	_ "ecom/cmd/docs"
	"ecom/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// @Summary      User login
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        login body ReqLogin true "Login credentials"
// @Success      200 {object} Resp "Successful login"
// @Failure      400 {object} util.ErrorResp "Invalid input"
// @Failure      401 {object} util.ErrorResp "Unauthorized"
// @Router       /login [post]
func login(ctx *fiber.Ctx) error {
	req := ReqLogin{}
	if err := ctx.BodyParser(&req); err != nil {
		return util.ErrorResponse(ctx, err)
	}

	if err := req.Validate(); err != nil {
		return util.ErrorResponse(ctx, err)
	}

	otpCode := util.GenerateRandomString(5)
	err := AppSrv.sms.Send(req.Phone, map[string]string{"code: ": otpCode})
	if err != nil {
		return util.ErrorResponse(ctx, err)
	}

	return nil
}

func otp(ctx *fiber.Ctx) error {
	return nil
}

func password(ctx *fiber.Ctx) error {
	return nil
}
