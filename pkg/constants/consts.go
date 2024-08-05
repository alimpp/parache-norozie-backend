package constants

import "errors"

// DebugMode is activated by passing debug=true argument.
// In DebugMode, endpoints don't check for jwt token.
// Be very careful not to activate DebugMode in production.
var DebugMode bool

const ServiceName = "Customer-Info"
const ServiceVersion = "2.0.0"
const MongoDbName = "customer_info"
const TemplateCollection = "template"
const AccountingCollection = "accounting"
const GoVersion = "1.22"
const CantParseJson = "cannot parse JSON"
const ErrAccountNotActivated = "account not activated"
const OutOfToken = "out of Token"
const NextFollowUpDate = "next_follow_up_date"
const DefaultFormsOrgId = "48e1c45e-2803-44ee-a4fb-6055ae92e9d6"

var (
	RecordNotFound = errors.New("record not found")
)

type ErrInvalidateData struct {
	ErrMsg string
}

func (e ErrInvalidateData) Error() string {
	return e.ErrMsg
}

type ErrNatsCon struct {
	ErrMsg string
}

func (e ErrNatsCon) Error() string {
	return e.ErrMsg
}
