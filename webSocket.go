package main

import(
	"encoding/json"
	"log"
	"net/http"

	"github.com/myselfBZ/Chat/models"
    
)

func HandleConn(w http.ResponseWriter, r *http.Request){
    var username map[string]string
    err := json.NewDecoder(r.Body).Decode(&username)
    if err != nil{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"err":"Enter your username",}) 
        return 

    }

    conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal("error upgrading the Conn:", err)
	}
    client := Client{
        Username: username["username"],
        Conn: conn,
    }
	Conns[client] = true 
	defer conn.Close() 
	
	for{
		var msg models.Message

		err := conn.ReadJSON(&msg)

		if err != nil {
			delete(Conns, client)
			break
		}		
		broadcast <- msg 

	
	}
}



func HandleMesg(){

	for{
		
		msg := <- broadcast

		for client := range Conns{
            if client.Username == msg.Reciever{
                err := client.Conn.WriteJSON(msg)
                if err != nil{
                    log.Fatal("error writing json:", err)
                    client.Conn.Close()
                    delete(Conns, client)
                }
            }
		}
	}
}
