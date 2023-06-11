package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
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
		send: make(chan []byte),
		Name: name,
	}

	client.hub.register <- client

	return client, nil
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
	}()

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Println("error: ", err)
			}
			break
		}

		log.Println("mensagem enviada: ", string(message))

		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)

			log.Println("mensagem recebida: ", string(message))

			if err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	client, err := newClient(hub, conn, r)

	if err != nil {
		log.Println(err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return
	}

	go client.readPump()
	go client.writePump()
}
