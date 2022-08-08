package competition

type Competition struct {
	// name of this Identifier
	name string
	// slug for this Identifier
	key string
}

func NewCompetition(name string, key string) Competition {
	return Competition{
		name: name,
		key:  key,
	}
}

func (c Competition) Key() string {
	return c.key
}

func (c Competition) Name() string {
	return c.name
}
