package constants

const ServiceName = "parche-norozie-backend"

type ErrRequestBody struct {
	Msg string
}

func (receiver ErrRequestBody) Error() string {
	return receiver.Msg
}
