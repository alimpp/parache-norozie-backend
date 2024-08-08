package constants

const ServiceName = "parche-norozie-backend"

type ErrRequestBody struct {
	Msg string
}

func (e ErrRequestBody) Error() string {
	return e.Msg
}

type ErrRecordNotFound struct {
	Msg string
}

func (e ErrRecordNotFound) Error() string {
	return e.Msg
}
