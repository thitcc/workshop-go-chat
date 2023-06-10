package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleChat(hub *Hub, w http.ResponseWriter, r *http.Request) {
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

	go client.handleMessage()
}

func main() {
	hub := newHub()
	go hub.run()

	log.Println("Servidor iniciado na porta 8080")

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handleChat(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
