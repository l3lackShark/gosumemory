package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

//JSONByte contains data that will be sent to the client
var JSONByte []byte

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	for {
		if memory.DynamicAddresses.IsReady == true {
			ws.WriteMessage(1, []byte(JSONByte)) //sending data to the client

		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}

//SetupRoutes creates websocket connection
func SetupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

//HTTPServer handles json and static files output
func HTTPServer() {

	for memory.DynamicAddresses.IsReady != true {
		time.Sleep(100 * time.Millisecond)
	}
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.Handle("/Songs/", http.StripPrefix("/Songs/", http.FileServer(http.Dir(memory.SongsFolderPath))))
	http.HandleFunc("/json", handler)
	http.ListenAndServe("127.0.0.1:24050", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	if memory.DynamicAddresses.IsReady == true {
		fmt.Fprintf(w, cast.ToString(JSONByte))

	} else {
		fmt.Fprintf(w, "osu! is not fully loaded!")
	}

}
