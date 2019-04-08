package main

import (
	"log"
	"sync"
	"time"
)

type PlayerState struct {
	ID string
	X  int
	Y  int
}

type ObjectState struct {
	ID   string
	Type string
	X    int
	Y    int
}

type RoomState struct {
	Players     []PlayerState
	Objects     []ObjectState
	CurrentTime time.Time
}

type Room struct {
	ID         string
	MaxPlayers uint
	Players    map[string]*Player
	mu         *sync.Mutex
	register   chan *Player
	unregister chan *Player
	ticker     *time.Ticker
	state      *RoomState
}

func NewRoom(maxPlayers uint) *Room {
	return &Room{
		MaxPlayers: maxPlayers,
		Players:    make(map[string]*Player),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		ticker:     time.NewTicker(1 * time.Second),
	}
}

func (r *Room) Run() {
	log.Println("room loop started")
	for {
		select {
		case player := <-r.unregister:
			delete(r.Players, player.ID)
			log.Println("player %s was removed from room", player.ID)
		case player := <-r.register:
			r.Players[player.ID] = player
			log.Printf("player %s joined", player.ID)
			player.SendMessage(&Message{"CONNECTED", nil})
		case <-r.ticker.C:
			log.Println("tick")

			// тут ваша игровая механика
			// взять команды у плеера, обработать их

			for _, player := range r.Players {
				player.SendState(r.state)
			}
		}
	}
}

func (r *Room) AddPlayer(player *Player) {
	player.room = r
	r.register <- player
}

func (r *Room) RemovePlayer(player *Player) {
	r.unregister <- player
}
