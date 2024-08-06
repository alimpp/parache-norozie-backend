package pkg

import (
	"ecom/pkg/constants"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(ctx *fiber.Ctx, err error) error {
	constants.Logger.Err(err).Msgf(err.Error())

	errorResp := errorResponse{Message: err.Error()}

	switch {
	case errors.Is(err, fiber.ErrUnprocessableEntity):
		errorResp.Message = "invalid body request"
		errorResp.Code = fiber.StatusUnprocessableEntity

	case errors.Is(err, constants.ErrRequestBody{}):
		errorResp.Message = "invalid body request"
		errorResp.Code = fiber.StatusBadRequest

	}

	return ctx.Status(fiber.StatusOK).JSON(errorResp)
}

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e errorResponse) Error() string {
	return e.Message
}
