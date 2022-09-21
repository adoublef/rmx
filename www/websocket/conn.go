package websocket

import (
	"github.com/gorilla/websocket"
	rmx "github.com/rog-golang-buddies/rapidmidiex/internal"
	"github.com/rog-golang-buddies/rapidmidiex/internal/suid"
)

type Conn struct {
	ID suid.UUID

	rwc *websocket.Conn
	p   *Pool
}

func (c Conn) Pool() *Pool { return c.p }

func (c Conn) Close() error {
	c.p.Delete(c.ID)

	return c.rwc.Close()
}

func (c Conn) ReadJSON(v any) error { return c.rwc.ReadJSON(v) }

func (c Conn) WriteJSON(v any) error { return c.rwc.WriteJSON(v) }

func (c Conn) SendMessage(v any) error {
	c.p.msgs <- v
	return nil
}

func (c Conn) SendMessage2(typ rmx.MessageType, data any) error {
	v := struct {
		Typ rmx.MessageType
	}{}
	c.p.msgs <- v
	return nil
}
