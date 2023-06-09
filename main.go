package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleChat(w http.ResponseWriter, r *http.Request) {
	log.Println("Conectado na rota /chat")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
}

func main() {
	log.Println("Servidor iniciado na porta 8080")

	http.HandleFunc("/chat", handleChat)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
