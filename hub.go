package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			for c := range h.clients {
				c.conn.WriteMessage(websocket.TextMessage, message)
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
