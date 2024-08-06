package api

type ReqLogin struct {
	Phone   string `json:"phone" validate:"required"`
	BackUrl string `json:"back_url"`
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
