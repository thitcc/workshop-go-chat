package main

import (
	"encoding/json"
	"log"

	"github.com/valyala/fastjson"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Event
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) sendMessage(message Event) {
	msg, err := json.Marshal(message.Message)

	if err != nil {
		return
	}

	to := string(fastjson.GetBytes(msg, "to"))
	from := string(fastjson.GetBytes(msg, "from"))

	if to == global {
		for client := range h.clients {
			client.send <- message
		}
	} else {
		for client := range h.clients {
			if to == client.Name || from == client.Name {
				client.send <- message
			}
		}
	}
}

func (h *Hub) notifyAll(event Event) {
	for c := range h.clients {
		c.send <- event
	}
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			h.sendMessage(message)
		case client := <-h.register:
			h.notifyAll(Event{Type: "user_connected", Message: client.Name})

			h.clients[client] = true

			log.Println("Usuário conectado: ", client.Name)
		case client := <-h.unregister:
			h.notifyAll(Event{Type: "user_disconnected", Message: client.Name})

			log.Println("Usuário desconectado: ", client.Name)

			client.conn.Close()
			delete(h.clients, client)
		}
	}
}
