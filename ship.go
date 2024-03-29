package battleship

type Ship struct {
	Type string
	Length int
	strength int
}

func (s *Ship) OnImpact() {
	s.strength -= 1
}

func (s *Ship) IsAfloat() bool {
	return s.strength > 0
}

func NewCruiser() *Ship {
	return &Ship{"Cruiser", 3, 3}
}

func NewSubmarine() *Ship {
	return &Ship{"Submarine", 3, 3}
}

func NewDestroyer() *Ship {
	return &Ship{"Destroyer", 2, 2}
}