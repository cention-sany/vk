package vk

const (
	ErrZero = iota
	ErrUnknown
	ErrAppDis
	ErrUnknownMethod
	ErrIncorrectSignature
	ErrAuthorizeFailed
	ErrTooManyReq
	ErrActDenied
	ErrInvalidReq
	ErrFloodControl
	ErrInternalServer
	ErrTestMode
	ErrTwelve
	ErrThirteen
	ErrCaptcha
	ErrAccDenied
	ErrHTTPAuthorizeFailed
	ErrValidationRequired
	ErrEighteen
	ErrNineteen
	ErrActDeniedNonStandalone
	ErrActOnlyStandalone
	ErrTwentyTwo
	ErrMethodDis
	ErrConfirmRequired
	ErrParamMissing     = 100
	ErrInvalidAppAPIID  = 101
	ErrInvalidUserID    = 113
	ErrInvalidTimestamp = 150
)

type Error struct {
	Code       int    `json:"error_code"`
	Msg        string `json:"error_msg"`
	CaptchaSId string `json:"captcha_sid"`
	CaptchaImg string `json:"captcha_img"`
	Redirect   string `json:"redirect_uri"`
}

// Implement error interface
func (e *Error) Error() string {
	return e.Msg
}
