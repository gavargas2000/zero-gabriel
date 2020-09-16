package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"encoding/json"
)

type child struct {
	Type string
	Timestamp int32
	Name string
}

type outputStruct struct {
	Type string
	Start int32
	End  int32
	Children []child
}

type inputItem struct {
	Timestamp float64 `json:"timestamp"`
	Type  string `json:"type"`
	SessionId  string `json:"session_id"`
	Name  string `json:"name"`
}

type inputStruct struct{
	items []inputItem
}

var currentSession = ""
var inputItems []inputItem

var ws_upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func read(conn *websocket.Conn) {
	var f interface{}

	for {
		_, byteStream, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}


		error := json.Unmarshal(byteStream, &f)



		if error != nil {
			//log.Printf("%v\n", f)

			//inputItems = append(inputItems, item)

			//	log.Printf("%v\n", inputStream)
		}

		for _, v := range f.([]interface{}) {
			var item inputItem

			tmpMap := v.(map[string]interface{})

			name, exists := tmpMap["name"]
			if exists {
				item.Name = name.(string)
			}

			id, exists := tmpMap["name"]
			if exists {
				item.SessionId = id.(string)
			}

			item.Type = tmpMap["type"].(string)
			item.Timestamp = tmpMap["timestamp"].(float64)

			inputItems = append(inputItems, item)
		}

		log.Println(fmt.Printf("%v\n", inputItems))
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := ws_upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	if err != nil {
		log.Println(err)
	}
	read(ws)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gabriel Vargas Zero Assignment")
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/ws", wsEndpoint)

	log.Fatal(http.ListenAndServe(":8844", router))

}





