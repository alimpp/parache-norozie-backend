package api

import (
	"crypto/subtle"
	_ "ecom/docs"
	"ecom/pkg/constants"
	"ecom/pkg/services"
	"ecom/pkg/util"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// @Summary      User otp login
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        login body LoginReq true "Login credentials"
// @Success      200 {object} Resp{data=LoginResp}
// @Router       /login [post]
func login(ctx *fiber.Ctx) error {
	req := LoginReq{}
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, err)
	}
	if err := req.Validate(); err != nil {
		return ErrorResponse(ctx, err)
	}

	err := AppSrv.redis.Get(ctx.Context(), req.Phone).Err()

	if err == nil {
		return ErrorResponse(ctx, constants.ErrDuplicateOtp{Msg: "برای درخواست مجدد صبر کنید"})
	}

	otpCode := util.GenerateRandomString(6)

	if err := AppSrv.redis.Set(ctx.Context(), req.Phone, otpCode, time.Duration(AppSrv.cfg.OTP.TTE)*time.Second).Err(); err != nil {
		return ErrorResponse(ctx, err)
	}

	if err := AppSrv.sms.Send(req.Phone, map[string]string{"code: ": otpCode}); err != nil {
		return ErrorResponse(ctx, err)
	}

	if AppSrv.debugMode {
		constants.Logger.Debug().Msgf("OTP code for Phone: %s is: %s", req.Phone, otpCode)
	}

	loginResp := LoginResp{TTE: AppSrv.cfg.OTP.TTE}

	if err := AppSrv.sqlDb.Where(services.User{Phone: req.Phone}).First(&services.User{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorResponse(ctx, err)
		}
	} else {
		loginResp.UserExists = true
	}

	resp := Resp{Data: loginResp, Status: 200}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary      Used to verify otp
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        OtpVerifyReq body OtpVerifyReq true "otp token"
// @Success      200 {object} Resp{data=OtpVerifyResp}
// @Router       /otp [post]
func verifyOtp(ctx *fiber.Ctx) error {
	req := OtpVerifyReq{}
	if err := ctx.BodyParser(&req); err != nil {
		fmt.Println(reflect.TypeOf(err))
		return ErrorResponse(ctx, err)
	}
	if err := req.Validate(); err != nil {
		return ErrorResponse(ctx, err)
	}

	strCmd := AppSrv.redis.Get(ctx.Context(), req.Phone)

	if strCmd.Err() != nil {
		if errors.Is(strCmd.Err(), redis.Nil) {
			return ErrorResponse(ctx, constants.ErrRecordNotFound{Msg: "کدی برای این شماره ارسال نشده"})
		} else {
			return ErrorResponse(ctx, strCmd.Err())
		}
	}

	byteCmd, _ := strCmd.Bytes()
	if equal := subtle.ConstantTimeCompare([]byte(req.OtpToken), byteCmd); equal != 1 {
		return ErrorResponse(ctx, constants.ErrBadRequest{Msg: "کد وارد شده اشتباه است"})
	}

	sess, err := AppSrv.sessionStore.Get(ctx)
	if err != nil {
		return ErrorResponse(ctx, err)
	}
	sess.Set("phone", req.Phone)

	if err := sess.Save(); err != nil {
		return ErrorResponse(ctx, err)
	}

	user := &services.User{}
	if err := AppSrv.sqlDb.FirstOrCreate(user, services.User{Phone: req.Phone}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorResponse(ctx, err)
		}
	}
	resp := Resp{Data: OtpVerifyResp{LoginSuccessful: true}, Status: 200}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary      Used to log out
// @Tags         authentication
// @Accept       plain
// @Produce      json
// @Success      200 {object} Resp
// @Router       /logout [post]
func logout(ctx *fiber.Ctx) error {
	sess, err := AppSrv.sessionStore.Get(ctx)
	if err != nil && sess != nil {
		AppSrv.sessionStore.Delete(sess.ID())
	}
	resp := Resp{Message: "از سرویس خارج شدید", Status: 200}
	return ctx.Status(fiber.StatusOK).JSON(resp)

}

func password(ctx *fiber.Ctx) error {
	return nil
}
