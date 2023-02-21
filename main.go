package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unibytes/fsgo/services"
)

func main() {
	r := mux.NewRouter()
	handler := services.NewFileHandler("static")
	r.HandleFunc("api/upload", handler.Upload).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}
