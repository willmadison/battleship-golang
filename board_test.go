package battleship

import "testing"
import "fmt"

func TestShipPlacementInvalidLocationRange(t *testing.T) {
	board := Board{4}
	cruiser := NewCruiser()
	r, _ := NewLocationRange("A4", "A1")
	err := board.Place(cruiser, r)

	if err == nil {
		t.Fatal("expected placement failure, got:", err)
	}
}

func TestLocationRange(t *testing.T) {
	board := Board{4}
	cruiser := NewCruiser()

	testCases := []struct {
		start, end  string
		validity    bool
		expectError bool
	}{
		{"A1", "A3", true, false},
		{"A3", "A1", false, true},
		{"A", "A1", false, true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("LocationRange(%s:%s)", testCase.start, testCase.end), func(t *testing.T) {
			r, err := NewLocationRange(testCase.start, testCase.end)

			if err != nil {
				if !testCase.expectError {
					t.Fatal("unexpected location range validation failure, got:", err)
				} else {
					return
				}
			} else if testCase.expectError {
				t.Fatal("expected location range validation failure but did not get any error")
			}

			valid, err := r.IsValidFor(board, cruiser)

			if valid != testCase.validity {
				t.Fatal("expected valid =", testCase.validity, "got:", valid,
					"for location range for", fmt.Sprintf("%s:%s", testCase.start, testCase.end))
			}

			if err != nil {
				if !testCase.expectError {
					t.Fatal("unexpected location range validation failure, got:", err)
				}
			} else if testCase.expectError {
				t.Fatal("expected location range validation failure but did not get any error")
			}
		})
	}
}
