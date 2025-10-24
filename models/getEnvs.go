package models

import (
	"diplomaGoSologub/tests"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func PortGetEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port string
	port = os.Getenv("TODO_PORT")
	// HTTP init
	//-- set port
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		fmt.Println("Error:", err)
		panic("port is not numeric")
	}

	return port
}

func WebDirGetEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webDir := os.Getenv("WEBDIR")

	if webDir == "" {
		fmt.Errorf("webDir path is empty")
		panic("webDir path is unavailiable")
	}

	return webDir
}

func DBGetEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = "scheduler.db" // название файла согласно условию задачи
		log.Println("dbFile path is empty. Default path ./scheduler.db will be used.")
	}

	return dbFile
}
