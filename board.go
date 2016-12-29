package battleship

type Board struct {
	Width int
}

func (b *Board) Place(ship *Ship, r LocationRange) error {
	return nil
}

type LocationRange struct {
	start, end string
}

func (l *LocationRange) IsValidFor(b Board, s *Ship) (bool, error) {
	return false, nil
}