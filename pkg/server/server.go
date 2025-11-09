package server

import (
	"diplomaGoSologub/models"
	"diplomaGoSologub/pkg/api"

	"log"
	"net/http"
)

func Start() {
	//--fileServer

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(models.WebDirGetEnv()))))

	api.Init()

	//--start Listen
	address := "0.0.0.0:" + models.PortGetEnv()
	log.Println("create link:", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
