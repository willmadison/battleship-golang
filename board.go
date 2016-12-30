package battleship

import (
	"errors"
	"strings"
	"strconv"
	"fmt"
	"unicode/utf8"
	"bytes"
)

type ShotResult string

const (
	Hit ShotResult = "H"
	Miss ShotResult = "M"
)

type Board struct {
	Width int
	occupants map[Location]*Ship
	shotsByLocation map[Location]ShotResult
}

func (b *Board) Place(ship *Ship, r LocationRange) error {
	if _, err := r.IsValidFor(*b, ship); err != nil {
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

func (b *Board) IsOccupied(l Location) bool {
	_, ok := b.occupants[l]
	return ok
}

func (b Board) String() string {
	return fmt.Sprintf("Board{Width:%d}", b.Width)
}

func (b Board) Display() string {
	rows := []string{}

	for row := 0; row < b.Width + 3; row++ {
		isHeaderRow := row == 0 || row == b.Width + 3 - 1

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

func (b *Board) Attack(l Location) string {
	var result string

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

type LocationRange struct {
	Length     int
	Locations  []Location
}

func (l *LocationRange) IsValidFor(b Board, s *Ship) (bool, error) {
	if l.Length > b.Width {
		return false, errors.New("range cannot be larger than the board")
	}

	if l.Length != s.Length {
		return false, errors.New("range must be equal to the length of the ship")
	}

	return true, nil
}

func NewLocationRange(start, end string) (LocationRange, error) {
	startLocation := Location(start)
	err := startLocation.Validate()
	if err != nil {
		return LocationRange{}, err
	}

	endLocation := Location(end)
	err = endLocation.Validate()
	if err != nil {
		return LocationRange{}, err
	}

	if !startLocation.IsBefore(endLocation) {
		return LocationRange{}, errors.New("invalid range, start location must come before the end location")
	}

	if startLocation.IsDiagonalTo(endLocation) {
		return LocationRange{}, errors.New("invalid range, start and end locations must not be diagonal to one another")
	}

	var length int

	if startLocation.InSameColumn(endLocation) {
		startingRune, _ := utf8.DecodeRuneInString(startLocation.Row())
		endingRune, _ := utf8.DecodeRuneInString(endLocation.Row())

		length = int(endingRune - startingRune) + 1
	} else if startLocation.InSameRow(endLocation) {
		length = endLocation.Column() - startLocation.Column() + 1
	}

	locations := []Location{startLocation}

	locations = append(locations, locationsInRange(startLocation, endLocation)...)
	locations = append(locations, endLocation)

	if length != len(locations) {
		return LocationRange{}, errors.New("invalid range, inappropriate number of locations for the determined length: length=" + strconv.Itoa(length) + " numLocations=" + strconv.Itoa(len(locations)) )
	}

	return LocationRange{length, locations}, nil
}

func locationsInRange(start, end Location) []Location {
	intermediates := []Location{}

	if start.InSameColumn(end) {
		startingRune, _ := utf8.DecodeRuneInString(start.Row())
		endingRune, _ := utf8.DecodeRuneInString(end.Row())

		for i := startingRune+1; i < endingRune; i++ {
			intermediates = append(intermediates, Location(fmt.Sprintf("%c%d", i, start.Column())))
		}
	} else if start.InSameRow(end) {
		for i := start.Column()+1; i < end.Column(); i++ {
			intermediates = append(intermediates, Location(fmt.Sprintf("%s%d", start.Row(), i)))
		}
	}

	return intermediates
}

type Location string

func (l Location) Row() string {
	return string(l[0:1])
}

func (l Location) Column() int {
	column, _ := strconv.Atoi(string(l[1:]))
	return column
}

func (l Location) Validate() error {
	if len(l) < 2 {
		return errors.New("location code must be at least 2 characters long")
	}

	row := string(l[0:1])

	if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", row) {
		return errors.New("location code must start with a letter")
	}

	column := l[1:]

	if _, err := strconv.Atoi(string(column)); err != nil {
		return err
	}

	return nil
}

func (l Location) IsBefore(other Location) bool {
	return l.Row() < other.Row() || l.Column() < other.Column()
}

func (l Location) IsDiagonalTo(other Location) bool {
	return l.Row() != other.Row() && l.Column() != other.Column()
}

func (l Location) InSameColumn(other Location) bool {
	return l.Column() == other.Column()
}

func (l Location) InSameRow(other Location) bool {
	return l.Row() == other.Row()
}
