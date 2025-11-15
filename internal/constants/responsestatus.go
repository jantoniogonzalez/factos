package constants

type ResponseStatus string

const (
	StatusSuccess  ResponseStatus = "success"
	StatusError    ResponseStatus = "error"
	StatusRedirect ResponseStatus = "redirect"
)
