package validator

func (v *Validator) ValidateUsername(field string) {
	var usernameField FieldName = "username"

	v.ValidateField(NotEmpty(field, usernameField))
	v.ValidateField(MaxCharacters(field, usernameField, 16))
	v.ValidateField(MinCharacters(field, usernameField, 4))
	v.ValidateField(NoSpecialCharacters(field, usernameField))
	v.ValidateField(NoTrailingSpecialCharacters(field, usernameField))
}
