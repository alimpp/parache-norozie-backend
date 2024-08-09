package api

import (
	"ecom/pkg/constants"
	"github.com/go-playground/validator/v10"
)

type Resp struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status"`
}

func (r Resp) Error() string {
	return r.Message
}

type LoginReq struct {
	Phone string `json:"phone" validate:"required"`
}

func (r LoginReq) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return constants.ErrRequestBody{Msg: "invalid json request body"}
	}
	return nil
}

type LoginResp struct {
	UserExists bool `json:"user_exists"`
	TTE        int  `json:"tte"`
}

//type ReqPassword struct {
//	Phone    string `json:"phone" validate:"required"`
//	Password string `json:"password" validate:"required"`
//}

type OtpVerifyReq struct {
	Phone    string `json:"phone" validate:"required"`
	OtpToken string `json:"otp_token" validate:"required"`
}

func (r OtpVerifyReq) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return constants.ErrRequestBody{Msg: "invalid json request body"}
	}
	return nil
}

type OtpVerifyResp struct {
	LoginSuccessful bool `json:"login_successful"`
}
