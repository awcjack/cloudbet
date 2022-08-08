package event

type Market struct {
	// mapping between submarket key and all associated submarkets for this Market
	submarkets map[string][]Selection
}

func NewMarket(submarkets map[string][]Selection) Market {
	return Market{
		submarkets: submarkets,
	}
}

func (m Market) UpdateSelection(key string, selections []Selection) {
	m.submarkets[key] = selections
}

func (m Market) Submarkets() map[string][]Selection {
	return m.submarkets
}
