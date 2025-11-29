package validator

import (
	"fmt"
	"strings"
	"unicode"
)

type ErrorMessage string
type FieldName string

type Validator struct {
	FieldErrors map[FieldName]ErrorMessage
}

// Return true if field does not have any field errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(fieldname FieldName, errMsg ErrorMessage) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[FieldName]ErrorMessage)
	}
	if _, exists := v.FieldErrors[fieldname]; !exists {
		v.FieldErrors[fieldname] = errMsg
	}
}

func (v *Validator) ValidateField(ok bool, fieldname FieldName, errMsg ErrorMessage) {
	if !ok {
		v.AddFieldError(fieldname, errMsg)
	}
}

// Returns true if it does not exceed max characters
func MaxCharacters(field string, fieldname FieldName, n int) (bool, FieldName, ErrorMessage) {
	return len(field) <= n,
		fieldname,
		ErrorMessage(fmt.Sprintf("This field cannot contain more than %d characters", n))
}

// Returns true if it contains the minimum amount of characters
func MinCharacters(field string, fieldname FieldName, n int) (bool, FieldName, ErrorMessage) {
	return len(field) >= n,
		fieldname,
		ErrorMessage(fmt.Sprintf("This field must contain at least %d characters", n))
}

// Returns true if is not empty
func NotEmpty(field string, fieldname FieldName) (bool, FieldName, ErrorMessage) {
	return strings.TrimSpace(field) != "",
		fieldname,
		ErrorMessage("This field cannot be empty")
}

// Returns true if it does not contain special characters
func NoSpecialCharacters(field string, fieldname FieldName) (bool, FieldName, ErrorMessage) {
	f := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' && r != '_'
	}

	return strings.IndexFunc(field, f) == -1,
		fieldname,
		ErrorMessage("This field must only contain alphanumeric character, period (.), and underscore (_)")
}

// Returns true if it does not start or finish with special characters
func NoTrailingSpecialCharacters(field string, fieldname FieldName) (bool, FieldName, ErrorMessage) {
	if field == "" {
		return false, fieldname, "This field cannot be empty"
	}

	runes := []rune(field)
	first := runes[0]
	last := runes[len(runes)-1]

	return (unicode.IsLetter(first) || unicode.IsNumber(first)) && (unicode.IsLetter(last) || unicode.IsNumber(last)),
		fieldname,
		ErrorMessage("This field cannot start or end with non-alphanumeric characters")
}

// Returns true if the field is not zero
func NotZero(field int, fieldname FieldName) (bool, FieldName, ErrorMessage) {
	return field != 0, fieldname, ErrorMessage("This field cannot be zero")
}
