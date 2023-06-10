package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleMessage(hub *Hub, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		log.Println("Usuário desconectado")
	}()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Println("error: ", err)
			}
			break
		}

		log.Println(string(message))

		hub.broadcast <- message
	}
}

func handleChat(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("Usuário conectado")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	hub.connections = append(hub.connections, conn)

	go handleMessage(hub, conn)
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
