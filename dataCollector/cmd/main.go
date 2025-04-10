package main

import (
	"dataCollector/cmd/server"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("must set path to ini file as only argument")
	}

	server.RunServer(os.Args[1])
}
