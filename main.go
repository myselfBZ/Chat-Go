package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/myselfBZ/Chat/models"
)

var(
	conns = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request)bool{
			return true
		},
	}
	broadcast = make(chan models.Message)
)


func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleConn)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "client/index.html")
	})
	go handleMesg()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Println("Server is running on http://localhost:8080/")
	server.ListenAndServe()
}


func handleConn(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal("error upgrading the conn:", err)
	}
	conns[conn] = true 
	defer conn.Close() 
	
	for{
		var msg models.Message

		err := conn.ReadJSON(&msg)

		if err != nil {
			delete(conns, conn)
			break
		}		
		broadcast <- msg 

	
	}
}

func handleMesg(){

	for{
		
		msg := <- broadcast
		log.Println(msg)
		for conn := range conns{
			err := conn.WriteJSON(msg)
			if err != nil{
				log.Fatal("error writing json:", err)
				conn.Close()
				delete(conns, conn)
			}
		}
	}
}











