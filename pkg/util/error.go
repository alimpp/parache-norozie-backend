package util

import (
	"ecom/pkg/constants"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(ctx *fiber.Ctx, err error) error {
	constants.Logger.Err(err).Msgf(err.Error())

	errorResp := ErrorResp{ErrMessage: err.Error()}

	switch {
	case errors.Is(err, fiber.ErrUnprocessableEntity):
		errorResp.ErrMessage = "invalid body request"
		errorResp.Status = fiber.StatusUnprocessableEntity

	case errors.Is(err, constants.ErrRequestBody{}):
		errorResp.ErrMessage = "invalid body request"
		errorResp.Status = fiber.StatusBadRequest
	default:
		errorResp.ErrMessage = "unknown error"
		errorResp.Status = fiber.StatusInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(errorResp)
}

type ErrorResp struct {
	ErrMessage string `json:"errMsg"`
	Status     int    `json:"status"`
}

func (e ErrorResp) Error() string {
	return e.ErrMessage
}
