package main


import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	Timestamp int32 `json:"timestamp"`
	Type  string `json:"type"`
	SessionId  string `json:"session_id"`
	Name  string `json:"name"`
}

type inputStruct struct{
	items []inputItem
}

var currentSession = ""




func inputEvents(w http.ResponseWriter, r *http.Request) {

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gabriel Vargas Zero Assignment")
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/input", inputEvents).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}





