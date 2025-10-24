package main

import (
	"diplomaGoSologub/models"
	"diplomaGoSologub/pkg/db"
	"diplomaGoSologub/pkg/server"
)

func main() {
	db.Init(models.DBGetEnv())

	server.Start()

}
