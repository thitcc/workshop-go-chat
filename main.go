package main

import (
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()

	log.Println("Servidor iniciado na porta 8080")

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
