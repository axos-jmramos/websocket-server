package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"web_socket/hub"

	"github.com/gorilla/websocket"
)

type AckResponse struct {
	Message string `json:"message"`
}

func main() {
	hub := hub.NewHub()
	id := 0
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//We need to authenticate the user first here.

		upgrader := websocket.Upgrader{}
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade", err)
			return
		}

		response := AckResponse{Message: "success"}
		c.WriteJSON(response)

		hub.Add(fmt.Sprintf("%v", id), c)
		log.Println("connection added")
		id++
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hub.NotifyUser("1", AckResponse{
			Message: "this is a test message",
		})
		w.Write([]byte(`{"status": "ok"}`))
		w.WriteHeader(http.StatusOK)
	})

	addr := flag.String("localhost", ":8000", "http service address")

	log.Print("Starting Webhook Server")
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
