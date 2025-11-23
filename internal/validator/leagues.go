package validator

import "github.com/jantoniogonzalez/factos/internal/models"

func (v *Validator) ValidateLeague(newLeague *models.League) {
	v.ValidateField(NotEmpty(newLeague.Name, "name"))
	v.ValidateField(NotEmpty(newLeague.Country, "country"))
	v.ValidateField(NotEmpty(newLeague.Logo, "logo"))
}
