package server

import (
	"diplomaGoSologub/models"
	"diplomaGoSologub/pkg/server/headers"

	"log"
	"net/http"
)

func Start() {
	//--fileServer
	http.HandleFunc("/test", headers.ServeHTTP)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(models.WebDirGetEnv()))))

	//--start Listen
	address := "127.0.0.1:" + models.PortGetEnv()
	log.Println(http.ListenAndServe(address, nil))
}
