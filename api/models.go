package api

import (
	"ecom/pkg/constants"
	"github.com/go-playground/validator/v10"
)

type ReqLogin struct {
	Phone   string `json:"phone" validate:"required"`
	BackUrl string `json:"back_url"`
}

func (r ReqLogin) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return constants.ErrRequestBody{Msg: err.Error()}
	}
	return nil
}

type ReqPassword struct {
	Phone    string `json:"phone" validate:"required"`
	BackUrl  string `json:"back_url"`
	Password string `json:"password" validate:"required"`
}

type Resp struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}
