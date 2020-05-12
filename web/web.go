package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

//JSONByte contains data that will be sent to the client
var JSONByte []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	if err != nil {
		fmt.Println("error:", err)
	}
	ws.WriteMessage(1, []byte(JSONByte)) //sending data to the client

}

//SetupRoutes creates websocket connection
func SetupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

//HTTPServer handles json and static files output
func HTTPServer() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/json", handler)
	http.ListenAndServe(":24050", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	if memory.DynamicAddresses.IsReady == true {
		fmt.Fprintf(w, cast.ToString(JSONByte))

	} else {
		fmt.Fprintf(w, "osu! is not fully loaded!")
	}

}
