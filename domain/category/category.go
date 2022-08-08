package category

type Category struct {
	// name of this Identifier
	name string
	// slug for this Identifier
	key string
}

func NewCategory(name string, key string) Category {
	return Category{
		name: name,
		key:  key,
	}
}

func (c Category) Key() string {
	return c.key
}

func (c Category) Name() string {
	return c.name
}
