package services

import "ecom/pkg/constants"

type SMS interface {
	Send(phone string, params map[string]string) error
}

type MockSMS struct {
}

func (s MockSMS) Send(phone string, params map[string]string) error {
	constants.Logger.Debug().Msgf("MockSMS server send phone number: %s, params: %s", phone, params)
	return nil
}
