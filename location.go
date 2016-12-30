package battleship

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// LocationRange represents a given span of continguous locations on a Battleship board.
type LocationRange struct {
	Length    int
	Locations []Location
}

// IsValidFor determines whether this LocationRange is valid for the given Board and Ship.
func (l *LocationRange) IsValidFor(b *Board, s *Ship) (bool, error) {
	if l.Length > b.Width {
		return false, errors.New("range cannot be larger than the board")
	}

	if l.Length != s.Length {
		return false, errors.New("range must be equal to the length of the ship")
	}

	return true, nil
}

// NewLocationRange returns a location range spanning the given start and end location codes.
// Returning an error if either location is invalid or creates an invalid range.
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

		length = int(endingRune-startingRune) + 1
	} else if startLocation.InSameRow(endLocation) {
		length = endLocation.Column() - startLocation.Column() + 1
	}

	locations := []Location{startLocation}

	locations = append(locations, locationsInRange(startLocation, endLocation)...)
	locations = append(locations, endLocation)

	if length != len(locations) {
		return LocationRange{}, errors.New("invalid range, inappropriate number of locations for the determined length: length=" + strconv.Itoa(length) + " numLocations=" + strconv.Itoa(len(locations)))
	}

	return LocationRange{length, locations}, nil
}

func locationsInRange(start, end Location) []Location {
	intermediates := []Location{}

	if start.InSameColumn(end) {
		startingRune, _ := utf8.DecodeRuneInString(start.Row())
		endingRune, _ := utf8.DecodeRuneInString(end.Row())

		for i := startingRune + 1; i < endingRune; i++ {
			intermediates = append(intermediates, Location(fmt.Sprintf("%c%d", i, start.Column())))
		}
	} else if start.InSameRow(end) {
		for i := start.Column() + 1; i < end.Column(); i++ {
			intermediates = append(intermediates, Location(fmt.Sprintf("%s%d", start.Row(), i)))
		}
	}

	return intermediates
}

// Location represents a discrete point on the Battleship board.
type Location string

// Row returns the row name as a string for this Location.
func (l Location) Row() string {
	return string(l[0:1])
}

// Column returns the Column numeral for this Location as an integer.
func (l Location) Column() int {
	column, _ := strconv.Atoi(string(l[1:]))
	return column
}

// Validate validates a given Location ensuring that it meets the requirements therein.
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

// IsBefore returns true if this Location is geographically before the given location, false otherwise.
func (l Location) IsBefore(other Location) bool {
	return l.Row() < other.Row() || l.Column() < other.Column()
}

// IsDiagonalTo returns true if this Location is geographically diagonal to the given location, false otherwise.
func (l Location) IsDiagonalTo(other Location) bool {
	return l.Row() != other.Row() && l.Column() != other.Column()
}

// InSameColumn returns true if this Location is geographically in the same column to the given location, false otherwise.
func (l Location) InSameColumn(other Location) bool {
	return l.Column() == other.Column()
}

// InSameRow returns true if this Location is geographically in the same row to the given location, false otherwise.
func (l Location) InSameRow(other Location) bool {
	return l.Row() == other.Row()
}
