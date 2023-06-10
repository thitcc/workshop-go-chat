package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var connections []*websocket.Conn

func handleMessage(conn *websocket.Conn) {
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

		for _, conn := range connections {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	log.Println("Usuário conectado")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	connections = append(connections, conn)

	go handleMessage(conn)
}

func main() {
	log.Println("Servidor iniciado na porta 8080")

	http.HandleFunc("/chat", handleChat)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
