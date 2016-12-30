package battleship

import "testing"
import (
	"fmt"
	"strconv"
)

func TestShipPlacementInvalidLocationRange(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A4", "A1")
	err := board.Place(cruiser, r)

	if err == nil {
		t.Fatal("expected placement failure, got:", err)
	}
}

func TestShipPlacementValidLocationRange(t *testing.T) {
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

func TestBoardDisplay(t *testing.T) {
	testCases := []struct {
		board Board
		want string
	}{
		{Board{Width: 4},
`===========
. 1 2 3 4
A
B
C
D
===========`},
		{Board{Width: 8},
`==================
. 1 2 3 4 5 6 7 8
A
B
C
D
E
F
G
H
==================`},
		{Board{Width: 12},
`=============================
. 1 2 3 4 5 6 7 8 9 10 11 12
A
B
C
D
E
F
G
H
I
J
K
L
=============================`},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Displaying %s", testCase.board), func(t *testing.T) {
			if testCase.want != testCase.board.Display() {
				t.Fatal("displayed:\n", testCase.board.Display(), "want:\n", testCase.want)
			}
		})
	}
}

func TestShipAttackSequence(t *testing.T) {
	board := Board{Width: 4}
	cruiser := NewCruiser()
	initialStrength := cruiser.strength
	r, _ := NewLocationRange("A1", "A3")
	err := board.Place(cruiser, r)

	if err != nil {
		t.Fatal("unexpected placement failure, got:", err)
	}

	result := board.Attack(Location("A2"))

	expected := "Hit. Cruiser."
	if result != expected {
		t.Fatal("incorrect attack result. expected " + expected + ", got: " + result)
	}

	if cruiser.strength != initialStrength - 1 {
		t.Fatal("incorrect attack result. cruiser strength to be at " + strconv.Itoa(initialStrength-1) + ", got: " + strconv.Itoa(cruiser.strength))
	}

	result = board.Attack(Location("A1"))

	result = board.Attack(Location("B1"))
	expected = "Miss!"
	if result != expected {
		t.Fatal("incorrect attack result. expected " + expected + ", got: " + result)
	}

	result = board.Attack(Location("A3"))

	expected = "Sunk Cruiser of length 3!"
	if result != expected {
		t.Fatal("incorrect attack result. expected " + expected + ", got: " + result)
	}

	expectedDisplay := `===========
. 1 2 3 4
A H H H
B M
C
D
===========`
	display := board.Display()

	if display != expectedDisplay {
		t.Fatal("displayed:\n", display, "want:\n", expectedDisplay)
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
