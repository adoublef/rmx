package rehearsal

import "github.com/rs/xid"

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

func NewRoom(s Title, c Capacity, o xid.ID) *Room {
	r := Room{xid.New(), s, c, o}
	return &r
}

type Title string

type Capacity int8
