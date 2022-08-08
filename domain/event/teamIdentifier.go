package event

type TeamIdentifier struct {
	// name of this Identifier
	name string
	// slug for this Identifier
	key string
	// abbreviation for this team's name
	abbreviation string
	// team country code
	nationality string
}

func NewTeamIdentifier(name string, key string, abbreviation string, nationality string) TeamIdentifier {
	if key == "" {
		return TeamIdentifier{}
	}

	return TeamIdentifier{
		name:         name,
		key:          key,
		abbreviation: abbreviation,
		nationality:  nationality,
	}
}

func (t TeamIdentifier) IsZero() bool {
	return t == TeamIdentifier{}
}

func (t TeamIdentifier) Name() string {
	return t.name
}

func (t TeamIdentifier) Key() string {
	return t.key
}

func (t TeamIdentifier) Abbreviation() string {
	return t.abbreviation
}

func (t TeamIdentifier) Nationality() string {
	return t.nationality
}
