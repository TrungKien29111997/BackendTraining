package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/demo", demoHanderler)
	log.Println("Server is Starting ...")
	err := http.ListenAndServe(":8080", nil) //local host
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
func demoHanderler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"message": "Hello, World!",
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
