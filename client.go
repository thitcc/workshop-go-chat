package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	Name string
}

func newClient(hub *Hub, conn *websocket.Conn, r *http.Request) (*Client, error) {
	params := r.URL.Query()

	name := params.Get("name")

	for c := range hub.clients {
		if name == c.Name {
			return nil, fmt.Errorf("esse nome já existe: %s", name)
		}
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		Name: name,
	}

	client.hub.register <- client

	return client, nil
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
