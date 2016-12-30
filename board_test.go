package battleship

import "testing"
import "fmt"

func TestShipPlacementInvalidLocationRange(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A4", "A1")
	err := board.Place(cruiser, r)

	if err == nil {
		t.Fatal("expected placement failure, got:", err)
	}
}

func TestShipPlacementValidLocation(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A1", "A3")
	err := board.Place(cruiser, r)

	if err != nil {
		t.Fatal("unexpected placement failure, got:", err)
	}
}

func TestShipPlacementOverlappingBoats(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A1", "A3")
	err := board.Place(cruiser, r)

	if err != nil {
		t.Fatal("unexpected placement failure, got:", err)
	}

	cruiser = NewCruiser()
	r, _ = NewLocationRange("A2", "C2")
	err = board.Place(cruiser, r)

	if err == nil {
		t.Fatal("expected placement failure, got:", err)
	}
}

func TestShipPlacementNonOverlappingBoats(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A1", "A3")
	err := board.Place(cruiser, r)

	if err != nil {
		t.Fatal("unexpected placement failure, got:", err)
	}

	cruiser = NewCruiser()
	r, _ = NewLocationRange("A4", "C4")
	err = board.Place(cruiser, r)

	if err != nil {
		t.Fatal("unexpected placement failure, got:", err)
	}

	for _, loc := range []Location{"A4", "B4", "C4"} {
		if !board.IsOccupied(loc) {
			t.Fatal("expected location " + string(loc) + " to be occupied")
		}
	}
}

func TestLocationRange(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()

	testCases := []struct {
		start, end  string
		expectedLength int
		validity    bool
		expectRangeError bool
		expectValidityError bool
	}{
		{"A1", "A3", 3, true, false, false},
		{"A1", "C1", 3, true, false, false},
		{"A3", "A1", 0, false, true, true},
		{"A", "A1", 0, false, true, true},
		{"AA", "A1", 0, false, true, true},
		{"A1", "C4", 0, false, true, true},
		{"A1", "A17", 17, false, false, true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("LocationRange(%s:%s)", testCase.start, testCase.end), func(t *testing.T) {
			r, err := NewLocationRange(testCase.start, testCase.end)

			if err != nil {
				if !testCase.expectRangeError {
					t.Fatal("unexpected location range validation failure, got:", err)
				} else {
					return
				}
			} else if testCase.expectRangeError {
				t.Fatal("expected location range validation failure but did not get any error")
			}

			if r.Length != testCase.expectedLength {
				t.Fatal("expected length of", testCase.expectedLength, "got:", r.Length)
			}

			valid, err := r.IsValidFor(board, cruiser)

			if valid != testCase.validity {
				t.Fatal("expected valid =", testCase.validity, "got:", valid,
					"for location range for", fmt.Sprintf("%s:%s", testCase.start, testCase.end))
			}

			if err != nil {
				if !testCase.expectValidityError {
					t.Fatal("unexpected location range validity failure, got:", err)
				}
			} else if testCase.expectValidityError {
				t.Fatal("expected location range validity failure but did not get any error")
			}
		})
	}
}
