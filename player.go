package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	conn *websocket.Conn
	ID   string
	in   chan *IncomeMessage
	out  chan *Message
	room *Room
}

func NewPlayer(conn *websocket.Conn, id string) *Player {
	return &Player{
		conn: conn,
		ID:   id,
		in:   make(chan *IncomeMessage),
		out:  make(chan *Message),
	}
}

func (p *Player) Listen() {
	go func() {
		for {
			message := &IncomeMessage{}
			err := p.conn.ReadJSON(message)
			if websocket.IsUnexpectedCloseError(err) {
				p.room.RemovePlayer(p)
				log.Println("player %s disconnected", p.ID)
				return
			}
			if err != nil {
				log.Printf("cannot read json")
				continue
			}

			p.in <- message
		}
	}()

	for {
		select {
		case message := <-p.out:
			p.conn.WriteJSON(message)
		case message := <-p.in:
			log.Printf("income: %#v", message)
		}
	}
}

func (p *Player) SendState(state *RoomState) {
	p.out <- &Message{"STATE", state}
}

func (p *Player) SendMessage(message *Message) {
	p.out <- message
}
