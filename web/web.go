package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/l3lackShark/gosumemory/config"
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
	if cast.ToBool(config.Config["cors"]) {
		enableCors(&w)
	}
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//HTTPServer handles json and static files output
func HTTPServer() {

	for memory.DynamicAddresses.IsReady != true {
		time.Sleep(100 * time.Millisecond)
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fs := http.FileServer(http.Dir(filepath.Join(exPath, "static")))
	http.Handle("/", fs)
	http.Handle("/Songs/", http.StripPrefix("/Songs/", http.FileServer(http.Dir(memory.SongsFolderPath))))
	http.HandleFunc("/json", handler)
	err = http.ListenAndServe(config.Config["serverip"], nil)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		log.Fatalln(err)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	if memory.DynamicAddresses.IsReady == true {
		fmt.Fprintf(w, cast.ToString(JSONByte))

	} else {
		fmt.Fprintf(w, `{"error": "osu! is not fully loaded!"}`)
	}

}
