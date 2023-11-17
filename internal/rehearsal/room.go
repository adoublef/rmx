package rehearsal

import (
	"errors"

	"github.com/rs/xid"
)

// Room for users to join and collaborate with
type Room struct {
	// ID is uniquely set for the room
	ID xid.ID
	// Title of the room
	Title Title
	// Capacity shares the maximum number of participants allowed
	Capacity Capacity
	// Owner is the id held by the user that created the room
	Owner xid.ID
}

// ParseRoom parses raw data into a Room construct
func ParseRoom(s string, n int, o xid.ID) (Room, error) {
	t, err := ParseTitle(s)
	if err != nil {
		return Room{}, err
	}
	c, err := ParseCapacity(n)
	if err != nil {
		return Room{}, err
	}
	return NewRoom(t, c, o), nil
}

// A New Room is constructed
func NewRoom(s Title, c Capacity, o xid.ID) Room {
	r := Room{xid.New(), s, c, o}
	return r
}

type Title string

// ParseTitle parses a string into the Title format.
func ParseTitle(t string) (Title, error) {
	if l := len(t); l < titleMin || l > titleMax {
		return "", errors.New("invalid length of title")
	}
	return Title(t), nil
}

const (
	titleMin = 1
	titleMax = 60
)

type Capacity int8

// ParseCapacity parses an integer into the Capacity format.
func ParseCapacity(c int) (Capacity, error) {
	if c < capMin || c > capMax {
		return 0, errors.New("invalid capacity size")
	}
	return Capacity(c), nil
}

const (
	capMin = 1
	capMax = 8
)
