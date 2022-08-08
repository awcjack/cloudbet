package sport

type Sport struct {
	// name of this Identifier
	name string
	// slug for this Identifier
	key string
	// average live time
	liveTime float64
}

func NewSport(name string, key string, liveTime float64) Sport {
	return Sport{
		name:     name,
		key:      key,
		liveTime: liveTime,
	}
}

func (s Sport) Key() string {
	return s.key
}

func (s Sport) Name() string {
	return s.name
}

func (s *Sport) SetLiveTime(time float64) {
	s.liveTime = time
}

func (s Sport) LiveTime() float64 {
	return s.liveTime
}
