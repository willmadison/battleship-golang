package battleship

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ShotResult represents either a successful or unsuccessful collision between an enemy missle and a vessel.
type ShotResult string

const (
	// Hit represents a successful collision between an enemy missle and a vessel.
	Hit ShotResult = "H"
	// Miss represents an unsuccessful collision between an enemy missle and a vessel.
	Miss ShotResult = "M"
)

// Board represents the notion of the Battleship battleground.
type Board struct {
	Width           int
	occupants       map[Location]*Ship
	shotsByLocation map[Location]ShotResult
}

// NewBoard returns a new battleship Board of the specified length.
func NewBoard(width int) *Board {
	board := &Board{Width: width}
	var ranges []LocationRange

	for i := 1; i <= width; i++ {
		rowName := fmt.Sprintf("%c", 'A'+i-1)
		start, end := fmt.Sprintf("%s%d", rowName, 1), fmt.Sprintf("%s%d", rowName, width)
		locationRange, err := NewLocationRange(start, end)
		if err != nil {
			fmt.Println("encountered an unexpected error creating a location range from", start, "to", end)
			continue
		}

		ranges = append(ranges, locationRange)
	}

	if board.occupants == nil {
		board.occupants = map[Location]*Ship{}
	}

	for _, r := range ranges {
		for _, l := range r.Locations {
			board.occupants[l] = nil
		}
	}

	return board
}

// Place places a given ship in the location range specified.
func (b *Board) Place(ship *Ship, r LocationRange) error {
	if _, err := r.IsValidFor(b, ship); err != nil {
		return err
	}

	for _, l := range r.Locations {
		if b.IsOccupied(l) {
			return errors.New("location " + string(l) + " is already occupied.")
		}
	}

	if b.occupants == nil {
		b.occupants = map[Location]*Ship{}
	}

	for _, l := range r.Locations {
		b.occupants[l] = ship
	}

	return nil
}

// IsOccupied returns true if the given location is occupied on the board, false otherwise.
func (b *Board) IsOccupied(l Location) bool {
	occupant, ok := b.occupants[l]
	return occupant != nil && ok
}

func (b Board) String() string {
	return fmt.Sprintf("Board{Width:%d}", b.Width)
}

// Display returns the Board's visual display including the status of each location therin as a string.
func (b Board) Display() string {
	rows := []string{}

	for row := 0; row < b.Width+3; row++ {
		isHeaderRow := row == 0 || row == b.Width+3-1

		if isHeaderRow {
			rows = append(rows, getHeader(b))
			continue
		}

		isColumnHeaderRow := row == 1

		if isColumnHeaderRow {
			rows = append(rows, getColumnHeader(b))
			continue
		}

		var derivedRow string
		var buffer bytes.Buffer

		for col := 0; col <= b.Width; col++ {
			isRowHeaderColumn := col == 0

			rowName := fmt.Sprintf(fmt.Sprintf("%c", 'A'+row-2))

			if isRowHeaderColumn {
				buffer.WriteString(rowName)
			} else {
				location := Location(fmt.Sprintf("%s%d", rowName, col))
				result, shotTaken := b.shotsByLocation[location]

				if shotTaken {
					buffer.WriteString(string(result))
				}
			}

			if col != b.Width {
				buffer.WriteString(" ")
			}

			derivedRow = strings.TrimSpace(buffer.String())
		}

		if len(derivedRow) > 0 {
			rows = append(rows, derivedRow)
		}
	}

	return strings.Join(rows, "\n")
}

func getHeader(b Board) string {
	var header string

	switch b.Width {
	case 4:
		header = "==========="
	case 8:
		header = "=================="
	case 12:
		header = "============================="
	}

	return header
}

func getColumnHeader(b Board) string {
	var buffer bytes.Buffer

	for i := 0; i <= b.Width; i++ {
		if i == 0 {
			buffer.WriteString(".")
		} else {
			buffer.WriteString(strconv.Itoa(i))
		}

		if i != b.Width {
			buffer.WriteString(" ")
		}
	}

	return buffer.String()
}

// Attack attempts an attack at the given location on the Board.
func (b *Board) Attack(l Location) string {
	var result string

	if !b.IsValidFor(l) {
		return fmt.Sprintf("Invalid Location %s. Please select another location.", string(l))
	}

	if b.shotsByLocation == nil {
		b.shotsByLocation = map[Location]ShotResult{}
	}

	if b.IsOccupied(l) {
		occupant := b.occupants[l]
		occupant.OnImpact()

		if occupant.IsAfloat() {
			result = fmt.Sprintf("Hit. %s.", occupant.Type)
		} else {
			result = fmt.Sprintf("Sunk %s of length %d!", occupant.Type, occupant.Length)
		}

		b.shotsByLocation[l] = Hit
	} else {
		result = "Miss!"
		b.shotsByLocation[l] = Miss
	}

	return result
}

// IsValidFor returns true if this given location is valid for this board, false otherwise.
func (b Board) IsValidFor(l Location) bool {
	_, ok := b.occupants[l]
	return ok
}
