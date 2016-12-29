package battleship

import (
	"errors"
	"strconv"
	"strings"
)

type Board struct {
	Width int
}

func (b *Board) Place(ship *Ship, r LocationRange) error {
	return nil
}

type LocationRange struct {
	start, end string
	Length     int
	locations  []Location
}

func NewLocationRange(start, end string) (LocationRange, error) {
	err := validateLocationCode(start)
	if err != nil {
		return LocationRange{}, err
	}

	err = validateLocationCode(end)
	if err != nil {
		return LocationRange{}, err
	}

	return LocationRange{}, nil
}

func validateLocationCode(code string) error {
	if len(code) < 2 {
		return errors.New("location code must be at least 2 characters long")
	}

	row := code[0:1]

	if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", row) {
		return errors.New("location code must start with a letter")
	}

	column := code[1:]

	if _, err := strconv.Atoi(column); err != nil {
		return err
	}

	return nil
}

type Location struct {
	Row, Column string
}

func (l *LocationRange) IsValidFor(b Board, s *Ship) (bool, error) {
	return false, nil
}
