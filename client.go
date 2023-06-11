package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
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
	send chan Event
	Name string
}

type Event struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
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
		send: make(chan Event),
		Name: name,
	}

	client.hub.register <- client

	return client, nil
}

func verifyMessage(message Event) bool {
	msg, _ := json.Marshal(message)

	return fastjson.Exists(msg, "type") &&
		fastjson.Exists(msg, "message") &&
		message.Type == "message"
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
		var message Event

		if err := c.conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Println("error: ", err)
			}
			break
		}

		if ok := verifyMessage(message); ok {
			c.hub.broadcast <- message
			log.Println("mensagem enviada: ", message)
		}
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

			if err := c.conn.WriteJSON(message); err != nil {
				log.Println(err)
			} else {
				log.Println("mensagem recebida: ", message)
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
