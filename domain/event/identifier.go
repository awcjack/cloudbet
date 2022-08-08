package event

import "errors"

type Identifier struct {
	// name of this Identifier
	name string
	// slug for this Identifier
	key string
}

func (i Identifier) Name() string {
	return i.name
}

func (i Identifier) Key() string {
	return i.key
}

func NewIdentifier(name string, key string) (Identifier, error) {
	if key == "" {
		return Identifier{}, errors.New("empty key")
	}

	return Identifier{
		name: name,
		key:  key,
	}, nil
}
