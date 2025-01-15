package validator

import (
	"strings"
)

type Validator struct {
	FieldErrors map[string]string
}

// Return true if field does not have any field errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(field, errMsg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = errMsg
	}
}

func (v *Validator) ValidateField(ok bool, field, errMsg string) {
	if !ok {
		v.AddFieldError(field, errMsg)
	}
}

// Returns true if it does not exceed max characters
func MaxCharacters(field string, n int) bool {
	return len(field) <= n
}

// Returns true if is not empty
func NotEmpty(field string) bool {
	return strings.TrimSpace(field) != ""
}
