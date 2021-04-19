package main

import (
	handler "GCP_image_upload/handler/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	httpLayer := handler.New() //DI
	r := mux.NewRouter()       // using gorilla mux router
	r.HandleFunc("/upload", httpLayer.Upload)
	log.Println("Server running in port:8080")
	_ = http.ListenAndServe(":8080", r)
}
