package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
}

func newClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{
		hub:  hub,
		conn: conn,
	}

	client.hub.register <- client

	return client
}

func (c *Client) handleMessage() {
	defer func() {
		c.hub.unregister <- c
	}()

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Println("error: ", err)
			}
			break
		}

		log.Println(string(message))

		c.hub.broadcast <- message
	}
}
