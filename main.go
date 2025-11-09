package main

import (
	"diplomaGoSologub/models"
	"diplomaGoSologub/pkg/db"
	"diplomaGoSologub/pkg/server"
	"log"
)

func main() {
	db.Init(models.DBGetEnv())

	log.Println("server start")
	server.Start()

}
