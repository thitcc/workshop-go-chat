package main

import (
	"log"
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

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		case client := <-h.register:
			h.clients[client] = true

			log.Println("Usuário conectado: ", client.Name)
		case client := <-h.unregister:
			client.conn.Close()
			delete(h.clients, client)

			log.Println("Usuário desconectado: ", client.Name)
		}
	}
}
