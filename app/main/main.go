package main

import (
	"log"
	"net/http"

	"github.com/Swetabh333/Makerble/app/routes"
)

func main() {
	router := routes.NewRouter()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Could not start server")
	}
}
