package event

type Selection struct {
	// outcome of this Selection
	outcome string
	// parameters to be sent by the client during bet placement on this selection, such as handicap, period etc.
	params string
	// price at which bets can be placed on this Selection
	price float64
	// maximum stake in EUR which can be placed in bets on this Selection; market liability = selection max stake * (price - 1); minimum stake is 0.01 EUR for all markets
	maxStake float64
	// probability of this Selection's outcome
	probability float64
	// current status of this Selection
	status string
	// side of this Selection (back/lay)
	side string
}

func NewSelection(outcome string, params string, price float64, maxStake float64, probability float64, status string, side string) Selection {
	return Selection{
		outcome:     outcome,
		params:      params,
		price:       price,
		maxStake:    maxStake,
		probability: probability,
		status:      status,
		side:        side,
	}
}

func (s Selection) Outcome() string {
	return s.outcome
}

func (s Selection) Params() string {
	return s.params
}

func (s Selection) Price() float64 {
	return s.price
}

func (s Selection) MaxStake() float64 {
	return s.maxStake
}

func (s Selection) Probability() float64 {
	return s.probability
}

func (s Selection) Status() string {
	return s.status
}

func (s Selection) Side() string {
	return s.side
}
