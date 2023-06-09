package main

import (
	"log"
	"net/http"
)

func handleChat(w http.ResponseWriter, r *http.Request) {
	log.Println("Conectado na rota /chat")
}

func main() {
	log.Println("Servidor iniciado na porta 8080")

	http.HandleFunc("/chat", handleChat)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
