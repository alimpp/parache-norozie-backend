package constants

const ServiceName = "parche_go"

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

type ErrDuplicateOtp struct {
	Msg string
}

func (e ErrDuplicateOtp) Error() string {
	return e.Msg
}

type ErrBadRequest struct {
	Msg string
}

func (e ErrBadRequest) Error() string {
	return e.Msg
}
