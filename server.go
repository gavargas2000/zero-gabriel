package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"encoding/json"
	"io/ioutil"
)

type child struct {
	Type string
	Timestamp int64
	Name string
}

type outputStruct struct {
	Type string
	Start int64
	End  int64
	Children []child
}

type inputItem struct {
	Timestamp int64 `json:"timestamp"`
	Type  string `json:"type"`
	SessionId  string `json:"session_id"`
	Name  string `json:"name"`
}

var currentSession = ""
var inputItems []inputItem
var outputItems outputStruct

var ws_upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func createFileName(name string) string{
	return fmt.Sprintf("%s%s", name, ".zero")
}

func storeItems(items outputStruct){
	var fileName = createFileName(currentSession)

	marshalledContent, _ := json.MarshalIndent(items, "", "\t")

	err := ioutil.WriteFile(fileName, marshalledContent, 0644)
	if err != nil {
		panic(err)
	}
}

func transformInput(items []inputItem){
	for _, item := range items {
		if item.Type == "SESSION_START"{
			outputItems.Start = item.Timestamp
			currentSession = item.SessionId
		}
		if item.Type == "SESSION_END"{
			outputItems.End = item.Timestamp
		}

		//SORT HERE?
		if item.Type == "EVENT" {
			var tmpChild child
			tmpChild.Type = item.Type
			tmpChild.Timestamp = item.Timestamp
			tmpChild.Name = item.Name

			outputItems.Children = append(outputItems.Children,tmpChild)
		}
	}
	outputItems.Type = "SESSION"
	storeItems(outputItems)
}

func read(conn *websocket.Conn) {
	var unmarshedInterface interface{}

	for {
		_, byteStream, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		error := json.Unmarshal(byteStream, &unmarshedInterface)

		if error != nil {

		}

		for _, v := range unmarshedInterface.([]interface{}) {
			var item inputItem

			tmpMap := v.(map[string]interface{})

			name, exists := tmpMap["name"]
			if exists {
				item.Name = name.(string)
			}

			id, exists := tmpMap["session_id"]
			if exists {
				item.SessionId = id.(string)
			}

			item.Type = tmpMap["type"].(string)
			item.Timestamp = int64(tmpMap["timestamp"].(float64))

			inputItems = append(inputItems, item)
		}

		transformInput(inputItems)
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

func returnSession(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	fileId := vars["id"]

	data, err := ioutil.ReadFile(createFileName(fileId))

	//Handle not found
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gabriel Vargas Zero Assignment")
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/websocket", wsEndpoint)
	router.HandleFunc("/session/{id}", returnSession).Methods("GET")

	log.Fatal(http.ListenAndServe(":8844", router))

}





