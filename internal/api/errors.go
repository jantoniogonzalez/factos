package api

import "errors"

var (
	ErrGoogleLoginFailed = errors.New("api: Google login failed")
	ErrReqMissingState   = errors.New("api: Request missing state")
	ErrStatesMismatch    = errors.New("api: States don't match")
	ErrMissingGoogleCode = errors.New("api: Missing Google code")
)
