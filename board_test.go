package battleship

import "testing"

func TestShipPlacementInvalidLocationRange(t *testing.T) {
	board := Board{4}
	cruiser := NewCruiser()
	r := LocationRange{"A4", "A1"}
	err := board.Place(cruiser, r)

	if err == nil {
		t.Fatal("expected placement failure, got:", err)
	}
}

func TestLocationRange(t *testing.T) {
	board := Board{4}
	cruiser := NewCruiser()
	r := LocationRange{"A4", "A1"}

	valid, err := r.IsValidFor(board, cruiser)

	if valid {
		t.Fatal("expected invalid location range for A4:A1")
	}

	if err == nil {
		t.Fatal("expected location range validation failure, got:", err)
	}
}
