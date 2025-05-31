package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Parameter Invalid
	ErrcodeTokenInvalid = 30001 // Token Invalid
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Parameter Invalid",
	ErrcodeTokenInvalid: "Token Invalid",
}
