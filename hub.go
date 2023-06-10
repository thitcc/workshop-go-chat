package main

import "github.com/gorilla/websocket"

type Hub struct {
	connections []*websocket.Conn
	broadcast   chan []byte
}

func newHub() *Hub {
	return &Hub{
		connections: []*websocket.Conn{},
		broadcast:   make(chan []byte),
	}
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			for _, conn := range h.connections {
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}
