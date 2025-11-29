package validator

func (v *Validator) ValidateLeagueName(name string) {
	v.ValidateField(NotEmpty(name, "name"))
}

func (v *Validator) ValidateLeagueCountry(country string) {
	v.ValidateField(NotEmpty(country, "country"))
}

func (v *Validator) ValidateLeagueLogo(logo string) {
	v.ValidateField(NotEmpty(logo, "logo"))
}

func (v *Validator) ValidateLeagueApiLeagueId(apiLeagueId int) {
	v.ValidateField(NotZero(apiLeagueId, "apiLeagueId"))
}

func (v *Validator) ValidateLeagueSeason(season int) {
	v.ValidateField(NotZero(season, "season"))
}
