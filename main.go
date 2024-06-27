package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/myselfBZ/Chat/middleware"

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
	mux.HandleFunc("/ws", middleware.AuthMiddleware(HandleConn))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "client/index.html")
	})
    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "client/login.html")
    })
    mux.HandleFunc("/log-in", func(w http.ResponseWriter, r *http.Request){
        var username struct{
            Username string `json:"username"`
        }
        if err := json.NewDecoder(r.Body).Decode(&username); err != nil{
            w.WriteHeader(http.StatusBadRequest)
            return 
        }
        token, err := middleware.GenerateToken(username.Username)
        if err != nil{
            w.WriteHeader(http.StatusBadRequest)
            return 
        }
        json.NewEncoder(w).Encode(map[string]string{"token":token})
    })
	go HandleMesg()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Println("Server is running on http://localhost:8080/")
	server.ListenAndServe()
}













