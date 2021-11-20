package main

import (
	"log"

	"github.com/doublen987/Projects/MySite/server/functionality"
	"github.com/doublen987/Projects/MySite/server/webportal"
)

type configuration struct {
	ServerAddress      string `json:"webserver"`
	DatabaseType       uint8  `json:"databasetype"`
	DatabaseConnection string `jdon:"dbconnection"`
	FrontEnd           string `json:"frontend"`
}

func main() {
	config, err := functionality.ExtractConfiguration("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	address := config.Host + ":" + config.Port
	log.Println("Starting web server on addres: ", address)

	err = webportal.RunAPI(config.Databasetype, address, config.DBConnection, config.FileStorageType)
	if err != nil {
		log.Fatal(err)
	}
}
