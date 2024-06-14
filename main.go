package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/myselfBZ/Chat/models"
)


type Client struct{
    Username string 
    Conn *websocket.Conn

}

var(
	Conns = make(map[Client]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request)bool{
			return true
		},
	}
	broadcast = make(chan models.Message)
)


func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", HandleConn)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "client/index.html")
	})
	go HandleMesg()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Println("Server is running on http://localhost:8080/")
	server.ListenAndServe()
}













