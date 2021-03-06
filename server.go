package main

/*
 * I think this file is still fairly small, so decided to keep it all in one file,
 * of course in a real life scenario you would probably have models, routers, etc on separate files.
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"encoding/json"
	"io/ioutil"
	"sort"
)

//Main struct to store children
type child struct {
	Type string
	Timestamp int64
	Name string
}

//Main struct for formatted output data
type outputStruct struct {
	Type string
	Start int64
	End  int64
	Children []child
}

//temp stuct to store data as it comes.
type inputItem struct {
	Timestamp int64 `json:"timestamp"`
	Type  string `json:"type"`
	SessionId  string `json:"session_id"`
	Name  string `json:"name"`
}

type errorStruct struct {
	Code int
	Message string
}

//main "global" vars
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

//This is where we are persisting the items, implemented easily just in files
//In a real world scenario this would be more of a DB or something like that
//Also decided to write the files Json formatted already.
func storeItems(items outputStruct){
	var fileName = createFileName(currentSession)

	sort.SliceStable(items.Children, func(i, j int) bool {
		return items.Children[i].Timestamp < items.Children[j].Timestamp
	})

	marshalledContent, _ := json.MarshalIndent(items, "", "\t")

	err := ioutil.WriteFile(fileName, marshalledContent, 0644)
	if err != nil {
		panic(err)
	}
}

//transforms the input into our main formatted output struct
func transformInput(items []inputItem) bool{
	var shouldClose = false

	for _, item := range items {
		if item.Type == "SESSION_START"{
			outputItems.Start = item.Timestamp

			if currentSession == ""{
				currentSession = item.SessionId
			}
		}
		if item.Type == "SESSION_END"{
			outputItems.End = item.Timestamp
			shouldClose = true
		}

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
	return shouldClose
}

//Read contents from the web socket and prepare it to process
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

		/*This whole block I think should be possible to do in a simpler way
		 * although for some reason this is the only way I made it work,
		 * so, sacrificed a more elegant solution for a working one here.
		 */
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

		shouldClose := transformInput(inputItems)
		log.Println("Should Close ", shouldClose)
		inputItems = []inputItem{}

		//Closes session at the end
		if shouldClose {
			currentSession = ""
			outputItems = outputStruct{}
			conn.Close()
			log.Println("Connection Closed")
		}
	}
}

func formatError(w http.ResponseWriter, err error, code int){
	var errorData errorStruct

	errorData.Code = code
	errorData.Message = err.Error()


	log.Printf("An error accured: %v", err)
	w.WriteHeader(code)

	w.Header().Set("Content-Type", "application/json")

	marshalledContent, _ := json.MarshalIndent(errorData, "", "\t")
	w.Write(marshalledContent)

}

//**** Endpoints and Views ***** ///

//Websocket Endpoint that reads the data, processs it,  stores it
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := ws_upgrader.Upgrade(w, r, nil)
	vars := mux.Vars(r)
	currentSession = vars["id"]

	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	if err != nil {
		log.Println(err)
	}
	read(ws)
}

//2nd endpoint to retreive a Session in formatted output
func returnSession(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	fileId := vars["id"]

	//Just reading from the files we stored
	data, err := ioutil.ReadFile(createFileName(fileId))

	//Handle not found
	if err != nil {
		formatError(w, err, 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gabriel Vargas Zero Assignment")
}

//Main function and Router
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")

	//Main Endpoint to send data via websocket
	router.HandleFunc("/websocket/{id}", wsEndpoint)
	router.HandleFunc("/websocket/", wsEndpoint)

	//Second endpoint to retreive a session in formatted output
	router.HandleFunc("/session/{id}", returnSession).Methods("GET")

	log.Fatal(http.ListenAndServe(":8844", router))
}





